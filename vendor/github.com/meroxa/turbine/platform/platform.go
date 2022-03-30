package platform

import (
	"context"
	"errors"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/turbine"
)

type Turbine struct {
	client    *Client
	functions map[string]turbine.Function
	deploy    bool
	imageName string
	config    turbine.AppConfig
	secrets   map[string]string
}

func New(deploy bool, imageName string) Turbine {
	c, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}

	ac, err := turbine.ReadAppConfig()
	if err != nil {
		log.Fatalln(err)
	}
	return Turbine{
		client:    c,
		functions: make(map[string]turbine.Function),
		imageName: imageName,
		deploy:    deploy,
		config:    ac,
		secrets:   make(map[string]string),
	}
}

func (t *Turbine) findPipeline(ctx context.Context) error {
	p, err := t.client.GetPipelineByName(ctx, t.config.Pipeline)
	if err != nil {
		return err
	}
	log.Printf("pipeline: %q (%q)", p.Name, p.UUID)

	return nil
}

func (t *Turbine) createPipeline(ctx context.Context) error {
	var input *meroxa.CreatePipelineInput

	input = &meroxa.CreatePipelineInput{
		Name: t.config.Pipeline,
		Metadata: map[string]interface{}{
			"app":     t.config.Name,
			"turbine": true,
		},
	}

	p, err := t.client.CreatePipeline(ctx, input)
	if err != nil {
		return err
	}

	// Alternatively, if we want to hide pipeline information completely by not logging this out,
	// we could create the application directly in Turbine
	log.Printf("pipeline: %q (%q)", p.Name, p.UUID)

	return nil
}

func (t Turbine) Resources(name string) (turbine.Resource, error) {
	if !t.deploy {
		return Resource{}, nil
	}

	ctx := context.Background()

	// Make sure we only create pipeline once
	if ok := t.findPipeline(ctx); ok != nil {
		err := t.createPipeline(ctx)
		if err != nil {
			return nil, err
		}
	}

	cr, err := t.client.GetResourceByNameOrID(ctx, name)
	if err != nil {
		return nil, err
	}

	log.Printf("retrieved resource %s (%s)", cr.Name, cr.Type)

	return Resource{
		ID:     cr.ID,
		Name:   cr.Name,
		Type:   string(cr.Type),
		client: t.client,
		v:      t,
	}, nil
}

type Resource struct {
	ID     int
	UUID   uuid.UUID
	Name   string
	Type   string
	client meroxa.Client
	v      Turbine
}

func (r Resource) Records(collection string, cfg turbine.ResourceConfigs) (turbine.Records, error) {
	if r.client == nil {
		return turbine.Records{}, nil
	}

	// TODO: ideally this should be handled on the platform
	mapCfg := cfg.ToMap()
	switch r.Type {
	case "redshift", "postgres", "mysql": // JDBC
		mapCfg["transforms"] = "createKey,extractInt"
		mapCfg["transforms.createKey.fields"] = "id"
		mapCfg["transforms.createKey.type"] = "org.apache.kafka.connect.transforms.ValueToKey"
		mapCfg["transforms.extractInt.field"] = "id"
		mapCfg["transforms.extractInt.type"] = "org.apache.kafka.connect.transforms.ExtractField$Key"
	}

	ci := &meroxa.CreateConnectorInput{
		ResourceID:    r.ID,
		Configuration: mapCfg,
		Type:          meroxa.ConnectorTypeSource,
		Input:         collection,
		PipelineName:  r.v.config.Pipeline,
	}

	con, err := r.client.CreateConnector(context.Background(), ci)
	if err != nil {
		return turbine.Records{}, err
	}

	outStreams := con.Streams["output"].([]interface{})

	// Get first output stream
	out := outStreams[0].(string)

	log.Printf("created source connector to resource %s and write records to stream %s from collection %s", r.Name, out, collection)
	return turbine.Records{
		Stream: out,
	}, nil
}

func (r Resource) Write(rr turbine.Records, collection string, cfg turbine.ResourceConfigs) error {
	// bail if dryrun
	if r.client == nil {
		return nil
	}

	// TODO: ideally this should be handled on the platform
	mapCfg := cfg.ToMap()
	switch r.Type {
	case "redshift", "postgres", "mysql": // JDBC sink
		mapCfg["table.name.format"] = strings.ToLower(collection)
		mapCfg["pk.mode"] = "record_value"
		mapCfg["pk.fields"] = "id"
		if r.Type != "redshift" {
			mapCfg["insert.mode"] = "upsert"
		}
	case "s3":
		mapCfg["aws_s3_prefix"] = strings.ToLower(collection) + "/"
		mapCfg["value.converter"] = "org.apache.kafka.connect.json.JsonConverter"
		mapCfg["value.converter.schemas.enable"] = "true"
		mapCfg["format.output.type"] = "jsonl"
		mapCfg["format.output.envelope"] = "true"
	}

	ci := &meroxa.CreateConnectorInput{
		ResourceID:    r.ID,
		Configuration: mapCfg,
		Type:          meroxa.ConnectorTypeDestination,
		Input:         rr.Stream,
		PipelineName:  r.v.config.Pipeline,
	}

	_, err := r.client.CreateConnector(context.Background(), ci)
	if err != nil {
		return err
	}
	log.Printf("created destination connector to resource %s and write records from stream %s to collection %s", r.Name, rr.Stream, collection)
	return nil
}

func (t Turbine) Process(rr turbine.Records, fn turbine.Function) (turbine.Records, turbine.RecordsWithErrors) {
	// register function
	funcName := strings.ToLower(reflect.TypeOf(fn).Name())
	t.functions[funcName] = fn

	var out turbine.Records
	var outE turbine.RecordsWithErrors

	if t.deploy {
		// create the function
		cfi := &meroxa.CreateFunctionInput{
			InputStream: rr.Stream,
			Image:       t.imageName,
			EnvVars:     t.secrets,
			Args:        []string{funcName},
			Pipeline:    meroxa.PipelineIdentifier{Name: t.config.Pipeline},
		}

		log.Printf("creating function %s ...", funcName)
		fnOut, err := t.client.CreateFunction(context.Background(), cfi)
		if err != nil {
			log.Panicf("unable to create function; err: %s", err.Error())
		}
		log.Printf("function %s created (%s)", funcName, fnOut.UUID)
		out.Stream = fnOut.OutputStream
	} else {
		// Not deploying, so map input stream to output stream
		out = rr
	}

	return out, outE
}

func (t Turbine) TriggerFunction(name string, in []turbine.Record) ([]turbine.Record, []turbine.RecordWithError) {
	log.Printf("Triggered function %s", name)
	return nil, nil
}

func (t Turbine) GetFunction(name string) (turbine.Function, bool) {
	fn, ok := t.functions[name]
	return fn, ok
}

func (t Turbine) ListFunctions() []string {
	var funcNames []string
	for name := range t.functions {
		funcNames = append(funcNames, name)
	}

	return funcNames
}

// RegisterSecret pulls environment variables with the same name and ships them as Env Vars for functions
func (t Turbine) RegisterSecret(name string) error {
	val := os.Getenv(name)
	if val == "" {
		return errors.New("secret is invalid or not set")
	}

	t.secrets[name] = val
	return nil
}
