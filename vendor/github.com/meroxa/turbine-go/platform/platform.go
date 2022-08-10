package platform

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/turbine-go"
)

type Turbine struct {
	client    *Client
	functions map[string]turbine.Function
	resources map[string]turbine.Resource
	deploy    bool
	imageName string
	config    turbine.AppConfig
	secrets   map[string]string
	gitSha    string
}

var pipelineUUID string

func New(deploy bool, imageName, appName, gitSha string) Turbine {
	c, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}

	ac, err := turbine.ReadAppConfig(appName, "")
	if err != nil {
		log.Fatalln(err)
	}
	return Turbine{
		client:    c,
		functions: make(map[string]turbine.Function),
		resources: make(map[string]turbine.Resource),
		imageName: imageName,
		deploy:    deploy,
		config:    ac,
		secrets:   make(map[string]string),
		gitSha:    gitSha,
	}
}

func (t *Turbine) findPipeline(ctx context.Context) error {
	_, err := t.client.GetPipelineByName(ctx, t.config.Pipeline)
	return err
}

func (t *Turbine) createPipeline(ctx context.Context) error {
	input := &meroxa.CreatePipelineInput{
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
	pipelineUUID = p.UUID
	return nil
}

func (t *Turbine) createApplication(ctx context.Context) error {
	inputCreateApp := &meroxa.CreateApplicationInput{
		Name:     t.config.Name,
		Language: "golang",
		GitSha:   t.gitSha,
		Pipeline: meroxa.EntityIdentifier{UUID: null.StringFrom(pipelineUUID)},
	}
	_, err := t.client.CreateApplication(ctx, inputCreateApp)
	return err
}

func (t Turbine) Resources(name string) (turbine.Resource, error) {
	if !t.deploy {
		t.resources[name] = Resource{}
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

	resource, err := t.client.GetResourceByNameOrID(ctx, name)
	if err != nil {
		return nil, err
	}

	log.Printf("retrieved resource %s (%s)", resource.Name, resource.Type)

	u, _ := uuid.Parse(resource.UUID)
	return Resource{
		UUID:   u,
		Name:   resource.Name,
		Type:   string(resource.Type),
		client: t.client,
		v:      t,
	}, nil
}

type Resource struct {
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

	ci := &meroxa.CreateConnectorInput{
		ResourceName:  r.Name,
		Configuration: cfg.ToMap(),
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

func (r Resource) Write(rr turbine.Records, collection string) error {
	return r.WriteWithConfig(rr, collection, turbine.ResourceConfigs{})
}

func (r Resource) WriteWithConfig(rr turbine.Records, collection string, cfg turbine.ResourceConfigs) error {
	// bail if dryrun
	if r.client == nil {
		return nil
	}

	connectorConfig := cfg.ToMap()
	switch r.Type {
	case "redshift", "postgres", "mysql", "sqlserver": // JDBC sink
		connectorConfig["table.name.format"] = strings.ToLower(collection)
	case "mongodb":
		connectorConfig["collection"] = strings.ToLower(collection)
	case "s3":
		connectorConfig["aws_s3_prefix"] = strings.ToLower(collection) + "/"
	case "snowflakedb":
		r := regexp.MustCompile("^[a-zA-Z]{1}[a-zA-Z0-9_]*$")
		matched := r.MatchString(collection)
		if !matched {
			return fmt.Errorf("%q is an invalid Snowflake name - must start with "+
				"a letter and contain only letters, numbers, and underscores", collection)
		}
		connectorConfig["snowflake.topic2table.map"] =
			fmt.Sprintf("%s:%s", rr.Stream, collection)
	}

	ci := &meroxa.CreateConnectorInput{
		ResourceName:  r.Name,
		Configuration: connectorConfig,
		Type:          meroxa.ConnectorTypeDestination,
		Input:         rr.Stream,
		PipelineName:  r.v.config.Pipeline,
	}

	_, err := r.client.CreateConnector(context.Background(), ci)
	if err != nil {
		return err
	}
	log.Printf("created destination connector to resource %s and write records from stream %s to collection %s", r.Name, rr.Stream, collection)

	err = r.v.createApplication(context.Background())
	if err != nil {
		return err
	}
	log.Printf("created application %q", r.v.config.Name)

	return nil
}

func (t Turbine) Process(rr turbine.Records, fn turbine.Function) turbine.Records {
	// register function
	funcName := strings.ToLower(reflect.TypeOf(fn).Name())
	t.functions[funcName] = fn

	var out turbine.Records

	if t.deploy {
		// create the function
		cfi := &meroxa.CreateFunctionInput{
			Name:        funcName,
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

	return out
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

func (t Turbine) ListResources() []string {
	var resourceNames []string
	for name := range t.resources {
		resourceNames = append(resourceNames, name)
	}

	return resourceNames
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
