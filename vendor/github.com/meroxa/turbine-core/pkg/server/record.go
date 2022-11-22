package server

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/ir"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pb.TurbineServiceServer = (*recordService)(nil)

type recordService struct {
	pb.UnimplementedTurbineServiceServer
	deploymentSpec ir.DeploymentSpec
	resources      []*pb.Resource
}

func NewRecordService() *recordService {
	return &recordService{}
}

func (s *recordService) Init(ctx context.Context, request *pb.InitRequest) (*emptypb.Empty, error) {
	s.deploymentSpec.Definition = ir.DefinitionSpec{
		GitSha: request.GetGitSHA(),
		Metadata: ir.MetadataSpec{
			Turbine: ir.TurbineSpec{
				Language: ir.Lang(strings.ToLower(request.GetLanguage().String())),
				Version:  request.GetTurbineVersion(),
			},
			SpecVersion: ir.LatestSpecVersion,
		},
	}

	return empty(), nil
}

func (s *recordService) GetResource(ctx context.Context, request *pb.GetResourceRequest) (*pb.Resource, error) {
	r := &pb.Resource{
		Name: request.GetName(),
	}

	s.resources = append(s.resources, r)
	return r, nil
}

func resourceConfigsToMap(configs []*pb.Config) map[string]interface{} {
	m := make(map[string]interface{})
	for _, rc := range configs {
		m[rc.GetField()] = rc.GetValue()
	}
	return m
}

func (s *recordService) ReadCollection(ctx context.Context, request *pb.ReadCollectionRequest) (*pb.Collection, error) {
	if request.GetCollection() == "" {
		return &pb.Collection{}, fmt.Errorf("please provide a collection name to 'read'")
	}

	for _, c := range s.deploymentSpec.Connectors {
		// Only one source per app allowed.
		if c.Type == ir.ConnectorSource {
			return &pb.Collection{}, fmt.Errorf("only one call to 'read' is allowed per Meroxa Data Application")
		}
	}

	s.deploymentSpec.Connectors = append(
		s.deploymentSpec.Connectors,
		ir.ConnectorSpec{
			Collection: request.GetCollection(),
			Resource:   request.GetResource().GetName(),
			Type:       ir.ConnectorSource,
			Config:     resourceConfigsToMap(request.GetConfigs().GetConfig()),
		},
	)

	return &pb.Collection{}, nil
}

func (s *recordService) WriteCollectionToResource(ctx context.Context, request *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	// This function may be called zero or more times.
	if request.GetTargetCollection() == "" {
		return empty(), fmt.Errorf("please provide a collection name to 'write'")
	}

	s.deploymentSpec.Connectors = append(
		s.deploymentSpec.Connectors,
		ir.ConnectorSpec{
			Collection: request.GetTargetCollection(),
			Resource:   request.GetResource().GetName(),
			Type:       ir.ConnectorDestination,
			Config:     resourceConfigsToMap(request.GetConfigs().GetConfig()),
		},
	)

	return empty(), nil
}

func (s *recordService) AddProcessToCollection(ctx context.Context, request *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	p := request.GetProcess()
	s.deploymentSpec.Functions = append(
		s.deploymentSpec.Functions,
		ir.FunctionSpec{
			Name: strings.ToLower(p.GetName()),
		})
	return &pb.Collection{}, nil
}

func (s *recordService) RegisterSecret(ctx context.Context, secret *pb.Secret) (*emptypb.Empty, error) {
	if s.deploymentSpec.Secrets == nil {
		s.deploymentSpec.Secrets = map[string]string{}
	}
	s.deploymentSpec.Secrets[secret.Name] = secret.Value
	return empty(), nil
}

func (s *recordService) HasFunctions(ctx context.Context, in *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	return wrapperspb.Bool(len(s.deploymentSpec.Functions) > 0), nil
}

func (s *recordService) ListResources(ctx context.Context, in *emptypb.Empty) (*pb.ListResourcesResponse, error) {
	return &pb.ListResourcesResponse{Resources: s.resources}, nil
}

func (s *recordService) GetSpec(ctx context.Context, in *pb.GetSpecRequest) (*pb.GetSpecResponse, error) {
	if image := in.GetImage(); image != "" {
		if len(s.deploymentSpec.Functions) == 0 {
			return nil, fmt.Errorf("cannot set function image since spec has no functions")
		}
		s.deploymentSpec.SetImageForFunctions(image)
	}

	spec, err := s.deploymentSpec.Marshal()
	if err != nil {
		return &pb.GetSpecResponse{}, err
	}

	return &pb.GetSpecResponse{
		Spec: spec,
	}, nil
}
