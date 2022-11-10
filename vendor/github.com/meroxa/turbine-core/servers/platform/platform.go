package platform

import (
	"context"
	"errors"
	"fmt"
	"log"

	turbinecore "github.com/meroxa/turbine-core"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/platform"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedTurbineServiceServer
	client       *platform.Client
	processes    map[string]pb.Process
	resources    []pb.Resource
	deploy       bool
	imageName    string
	config       turbinecore.AppConfig
	secrets      map[string]string
	gitSha       string
	appUUID      string
	pipelineUUID string
	PipelineName string
}

func (s *server) Init(ctx context.Context, request *pb.InitRequest) (*emptypb.Empty, error) {
	// create Meroxa Data Platform client
	c, err := platform.NewClient()
	if err != nil {
		return Empty(), err
	}
	s.client = c

	// load app config (app.json)
	ac, err := turbinecore.ReadAppConfig(request.GetAppName(), request.GetConfigFilePath())
	if err != nil {
		return Empty(), err
	}
	s.config = ac

	// Create the pipeline
	_, err = s.findOrCreatePipeline(ctx)
	if err != nil {
		return Empty(), err
	}
	log.Printf("pipeline: %s", s.pipelineUUID)

	// Create application
	_, err = s.findOrCreateApp(ctx, request.Language.String(), request.GitSHA)
	if err != nil {
		return Empty(), err
	}
	log.Printf("app: %s", s.appUUID)

	return Empty(), nil
}

func (s *server) GetResource(ctx context.Context, id *pb.NameOrUUID) (*pb.Resource, error) {
	var val string
	if uid := id.GetUuid(); uid != "" {
		val = uid
	} else if name := id.GetName(); name != "" {
		val = name
	} else {
		return nil, errors.New("name or UUID is required")
	}
	resource, err := s.client.GetResourceByNameOrID(ctx, val)
	if err != nil {
		return nil, err
	}

	log.Printf("retrieved resource %s (%s)", resource.Name, resource.Type)

	return &pb.Resource{
		Uuid: resource.UUID,
		Name: resource.Name,
		Type: string(resource.Type),
	}, nil
}

func resourceConfigsToMap(configs []*pb.ResourceConfig) map[string]interface{} {
	m := make(map[string]interface{})
	for _, rc := range configs {
		m[rc.GetField()] = rc.GetValue()
	}
	return m
}

func (s *server) ReadCollection(ctx context.Context, request *pb.ReadCollectionRequest) (*pb.Collection, error) {
	ci := &meroxa.CreateConnectorInput{
		ResourceName:  request.Resource.Name,
		Configuration: resourceConfigsToMap(request.GetConfigs().GetResourceConfig()),
		Type:          meroxa.ConnectorTypeSource,
		Input:         request.GetCollection(),
		PipelineName:  s.PipelineName,
	}

	con, err := s.client.CreateConnector(ctx, ci)
	if err != nil {
		return nil, err
	}

	outStreams := con.Streams["output"].([]interface{})

	// Get first output stream
	out := outStreams[0].(string)

	log.Printf("created source connector to resource %s and write records to stream %s from collection %s", request.Resource.Name, out, request.Collection)
	return &pb.Collection{
		Stream: out,
	}, nil
}

func (s *server) WriteCollectionToResource(ctx context.Context, request *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	log.Printf("request: %+v", request)
	var connectorConfig map[string]interface{}
	if request.Configs != nil {
		log.Printf("connectorConfig: %+v", connectorConfig)
		connectorConfig = resourceConfigsToMap(request.Configs.ResourceConfig)
	} else {
		connectorConfig = make(map[string]interface{})
	}

	ci := &meroxa.CreateConnectorInput{
		ResourceName:  request.Resource.Name,
		Configuration: connectorConfig,
		Type:          meroxa.ConnectorTypeDestination,
		Input:         request.Collection.Stream,
		PipelineName:  s.PipelineName,
	}

	_, err := s.client.CreateConnector(ctx, ci)
	if err != nil {
		return Empty(), err
	}
	log.Printf("created destination connector to resource %s and write records from stream %s to collection %s", request.Resource.Name, request.Collection.Stream, request.Collection)

	return Empty(), nil
}

func (s *server) AddProcessToCollection(ctx context.Context, request *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	funcNameGitSHA := fmt.Sprintf("%s-%.8s", request.Process.Name, s.gitSha)
	// create the function
	cfi := &meroxa.CreateFunctionInput{
		Name:        funcNameGitSHA,
		InputStream: request.Collection.Stream,
		Image:       request.Process.Image,
		EnvVars:     s.secrets,
		Args:        []string{request.Process.Name},
		Pipeline:    meroxa.PipelineIdentifier{Name: s.config.Pipeline},
	}

	log.Printf("creating function %s ...", request.Process.Name)
	fnOut, err := s.client.CreateFunction(ctx, cfi)
	if err != nil {
		log.Panicf("unable to create function; err: %s", err.Error())
	}
	log.Printf("function %s created (%s)", request.Process.Name, fnOut.UUID)
	return &pb.Collection{Stream: fnOut.OutputStream}, nil
}

func (s *server) RegisterSecret(ctx context.Context, secret *pb.Secret) (*emptypb.Empty, error) {
	s.secrets[secret.Name] = secret.Value
	return Empty(), nil
}

func (s *server) findOrCreatePipeline(ctx context.Context) (string, error) {
	if s.pipelineUUID != "" && s.PipelineName != "" {
		log.Printf("pipeline already persisted: %s", s.pipelineUUID)
		return s.pipelineUUID, nil
	}
	// todo: try finding pipeline first
	input := &meroxa.CreatePipelineInput{
		Name: s.config.Pipeline,
		Metadata: map[string]interface{}{
			"app":     s.config.Name,
			"turbine": true,
		},
	}

	p, err := s.client.CreatePipeline(ctx, input)
	if err != nil {
		log.Printf("unable to create pipeline: %s", err)
		return "", err
	}
	s.pipelineUUID = p.UUID
	s.PipelineName = p.Name
	return s.pipelineUUID, nil
}

func (s *server) findOrCreateApp(ctx context.Context, language string, gitSHA string) (string, error) {
	if s.appUUID != "" {
		log.Printf("app already persisted: %s", s.appUUID)
	}
	// todo: try finding app first
	inputCreateApp := &meroxa.CreateApplicationInput{
		Name:     s.config.Name,
		Language: language,
		GitSha:   gitSHA,
		Pipeline: meroxa.EntityIdentifier{Name: s.config.Pipeline},
	}

	a, err := s.client.CreateApplication(ctx, inputCreateApp)
	if err != nil {
		return "", err
	}
	s.appUUID = a.UUID
	return s.appUUID, nil
}

func New() *server {
	return &server{
		processes: make(map[string]pb.Process),
		resources: []pb.Resource{},
		deploy:    false,
		imageName: "",
		config:    turbinecore.AppConfig{},
		secrets:   nil,
		gitSha:    "",
		appUUID:   "",
	}
}

func Empty() *emptypb.Empty {
	return new(emptypb.Empty)
}
