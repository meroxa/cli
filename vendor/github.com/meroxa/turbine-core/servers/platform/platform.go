package platform

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/ir"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedTurbineServiceServer
	deploymentSpec ir.DeploymentSpec
	resources      []pb.Resource
}

func (s *server) Init(ctx context.Context, request *pb.InitRequest) (*emptypb.Empty, error) {
	s.deploymentSpec.Definition = ir.DefinitionSpec{
		GitSha: request.GitSHA,
		Metadata: ir.MetadataSpec{
			Turbine: ir.TurbineSpec{
				Language: ir.Lang(request.Language.String()),
				Version:  request.TurbineVersion,
			},
			SpecVersion: ir.LatestSpecVersion,
		},
	}

	return Empty(), nil
}

func (s *server) GetResource(ctx context.Context, id *pb.NameOrUUID) (*pb.Resource, error) {
	r := pb.Resource{
		Uuid: id.Uuid,
		Name: id.Name,
	}

	s.resources = append(s.resources, r)
	return &r, nil
}

func resourceConfigsToMap(configs []*pb.ResourceConfig) map[string]interface{} {
	m := make(map[string]interface{})
	for _, rc := range configs {
		m[rc.GetField()] = rc.GetValue()
	}
	return m
}

func (s *server) ReadCollection(ctx context.Context, request *pb.ReadCollectionRequest) (*pb.Collection, error) {
	if request.Collection == "" {
		return &pb.Collection{}, fmt.Errorf("please provide a collection name to Records()")
	}

	for _, c := range s.deploymentSpec.Connectors {
		// Only one source per app allowed.
		if c.Type == ir.ConnectorSource {
			return &pb.Collection{}, fmt.Errorf("only one call to Records() is allowed per Meroxa Data Application")
		}
	}

	s.deploymentSpec.Connectors = append(
		s.deploymentSpec.Connectors,
		ir.ConnectorSpec{
			Collection: request.Collection,
			Resource:   request.Resource.Name,
			Type:       ir.ConnectorSource,
			Config:     resourceConfigsToMap(request.GetConfigs().GetResourceConfig()),
		},
	)

	return &pb.Collection{}, nil
}

func (s *server) WriteCollectionToResource(ctx context.Context, request *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	// This function may be called zero or more times.
	if request.TargetCollection == "" {
		return Empty(), fmt.Errorf("please provide a collection name to WriteWithConfig()")
	}

	s.deploymentSpec.Connectors = append(
		s.deploymentSpec.Connectors,
		ir.ConnectorSpec{
			Collection: request.TargetCollection,
			Resource:   request.Resource.Name,
			Type:       ir.ConnectorDestination,
			Config:     resourceConfigsToMap(request.GetConfigs().GetResourceConfig()),
		},
	)

	fmt.Printf("DEPLOYMENT SPEC WriteCollectionToResource: %+#v", s.deploymentSpec)

	return Empty(), nil
}

func (s *server) AddProcessToCollection(ctx context.Context, request *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	p := request.Process
	s.deploymentSpec.Functions = append(
		s.deploymentSpec.Functions,
		ir.FunctionSpec{
			Name:  strings.ToLower(p.Name),
			Image: p.Image,
		})
	fmt.Printf("DEPLOYMENT SPEC AddProcessToCollection: %+#v", s.deploymentSpec)

	return nil, nil
}

func (s *server) RegisterSecret(ctx context.Context, secret *pb.Secret) (*emptypb.Empty, error) {
	s.deploymentSpec.Secrets[secret.Name] = secret.Value
	return Empty(), nil
}

func New() *server {
	return &server{}
}

func Empty() *emptypb.Empty {
	return new(emptypb.Empty)
}
