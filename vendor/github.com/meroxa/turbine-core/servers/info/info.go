package info

import (
	"context"
	"fmt"

	turbinecore "github.com/meroxa/turbine-core"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/platform"
	platform2 "github.com/meroxa/turbine-core/servers/platform"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedTurbineServiceServer
	client                   *platform.Client
	processes                map[string]*pb.Process
	resources                []pb.Resource
	resourcesWithCollections []pb.ResourceWithCollection
	deploy                   bool
	imageName                string
	config                   turbinecore.AppConfig
	secrets                  map[string]string
	gitSha                   string
	appUUID                  string
}

func (s *server) Init(ctx context.Context, request *pb.InitRequest) (*emptypb.Empty, error) {
	// load app config (app.json)
	ac, err := turbinecore.ReadAppConfig(request.GetAppName(), request.GetConfigFilePath())
	if err != nil {
		return platform2.Empty(), err
	}
	s.config = ac
	return platform2.Empty(), nil
}

func (s *server) ListFunctions(ctx context.Context, empty *emptypb.Empty) (*pb.ListFunctionsResponse, error) {
	ff := pb.ListFunctionsResponse{Functions: []string{}}
	for _, p := range s.processes {
		ff.Functions = append(ff.Functions, p.Name)
	}
	return &ff, nil
}

func (s *server) ListResources(ctx context.Context, empty *emptypb.Empty) (*pb.ListResourcesResponse, error) {
	rr := pb.ListResourcesResponse{Resources: []*pb.ResourceWithCollection{}}
	for _, r := range s.resources {
		var dir pb.ResourceWithCollection_Direction
		if r.Direction == pb.Resource_DIRECTION_SOURCE {
			dir = pb.ResourceWithCollection_DIRECTION_SOURCE
		} else {
			dir = pb.ResourceWithCollection_DIRECTION_DESTINATION
		}
		rr.Resources = append(rr.Resources, &pb.ResourceWithCollection{
			Name:       r.Name,
			Collection: "todo", //todo
			Direction:  dir,
		})
	}
	return &rr, nil
}

func (s *server) GetResource(ctx context.Context, id *pb.NameOrUUID) (*pb.Resource, error) {
	r := pb.Resource{}
	if id.Uuid != "" {
		r.Uuid = id.Uuid
	}
	if id.Name != "" {
		r.Name = id.Name
	}

	s.resources = append(s.resources, r)

	return &r, nil
}

func (s *server) ReadCollection(ctx context.Context, request *pb.ReadCollectionRequest) (*pb.Collection, error) {
	if request.Collection == "" {
		return &pb.Collection{}, fmt.Errorf("please provide a collection name to Records()")
	}

	s.resourcesWithCollections = append(s.resourcesWithCollections, pb.ResourceWithCollection{
		Name:       request.Resource.Name,
		Collection: request.Collection,
		Direction:  pb.ResourceWithCollection_DIRECTION_SOURCE,
	})

	return nil, nil
}

func (s *server) WriteCollectionToResource(ctx context.Context, request *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	if request.Collection.Stream == "" {
		return platform2.Empty(), fmt.Errorf("please provide a collection name to Records()")
	}

	s.resourcesWithCollections = append(s.resourcesWithCollections, pb.ResourceWithCollection{
		Name:       request.Resource.Name,
		Collection: request.Collection.Stream,
		Direction:  pb.ResourceWithCollection_DIRECTION_DESTINATION,
	})

	return platform2.Empty(), nil
}

func (s *server) AddProcessToCollection(ctx context.Context, request *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	s.processes[request.Process.Name] = request.Process
	return nil, nil
}

func New() *server {
	return &server{
		processes: make(map[string]*pb.Process),
		resources: []pb.Resource{},
		deploy:    false,
		imageName: "",
		config:    turbinecore.AppConfig{},
		secrets:   nil,
		gitSha:    "",
		appUUID:   "",
	}
}
