package server

import (
	"context"
	"strings"

	"github.com/google/uuid"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/ir"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pb.TurbineServiceServer = (*specBuilderService)(nil)

type specBuilderService struct {
	pb.UnimplementedTurbineServiceServer

	spec      *ir.DeploymentSpec
	resources []*pb.Resource
}

func NewSpecBuilderService() *specBuilderService {
	return &specBuilderService{
		spec: &ir.DeploymentSpec{
			Secrets: make(map[string]string),
		},
	}
}

func (s *specBuilderService) Init(_ context.Context, req *pb.InitRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.spec.Definition = ir.DefinitionSpec{
		GitSha: req.GetGitSHA(),
		Metadata: ir.MetadataSpec{
			Turbine: ir.TurbineSpec{
				Language: ir.Lang(strings.ToLower(req.Language.String())),
				Version:  req.TurbineVersion,
			},
			SpecVersion: ir.LatestSpecVersion,
		},
	}
	return empty(), nil
}

func (s *specBuilderService) GetResource(_ context.Context, req *pb.GetResourceRequest) (*pb.Resource, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &pb.Resource{Name: req.Name}, nil
}

func (s *specBuilderService) ReadCollection(_ context.Context, req *pb.ReadCollectionRequest) (*pb.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.resources = append(s.resources, &pb.Resource{
		Name:       req.GetResource().GetName(),
		Source:     true,
		Collection: req.GetCollection(),
	})

	c := ir.ConnectorSpec{
		UUID:       uuid.New().String(),
		Collection: req.Collection,
		Resource:   req.Resource.Name,
		Type:       ir.ConnectorSource,
		Config:     configMap(req.Configs),
	}

	if err := s.spec.AddSource(&c); err != nil {
		return nil, err
	}

	return &pb.Collection{
		Name:   req.Collection,
		Stream: c.UUID,
	}, nil
}

func (s *specBuilderService) WriteCollectionToResource(_ context.Context, req *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.resources = append(s.resources, &pb.Resource{
		Name:        req.Resource.Name,
		Destination: true,
		Collection:  req.TargetCollection,
	})

	c := ir.ConnectorSpec{
		UUID:       uuid.New().String(),
		Collection: req.TargetCollection,
		Resource:   req.Resource.Name,
		Type:       ir.ConnectorDestination,
		Config:     configMap(req.Configs),
	}
	if err := s.spec.AddDestination(&c); err != nil {
		return nil, err
	}

	if err := s.spec.AddStream(&ir.StreamSpec{
		UUID:     uuid.New().String(),
		FromUUID: req.SourceCollection.Stream,
		ToUUID:   c.UUID,
		Name:     req.SourceCollection.Stream + "_" + c.UUID,
	}); err != nil {
		return nil, err
	}

	return empty(), nil
}

func (s *specBuilderService) AddProcessToCollection(_ context.Context, req *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	f := ir.FunctionSpec{
		UUID: uuid.New().String(),
		Name: strings.ToLower(req.Process.Name),
		Script: req.Process.Script,
	}
	if err := s.spec.AddFunction(&f); err != nil {
		return nil, err
	}

	if err := s.spec.AddStream(&ir.StreamSpec{
		UUID:     uuid.New().String(),
		FromUUID: req.Collection.Stream,
		ToUUID:   f.UUID,
		Name:     req.Collection.Stream + "_" + f.UUID,
	}); err != nil {
		return nil, err
	}

	return &pb.Collection{
		Name:   req.Collection.Name,
		Stream: f.UUID,
	}, nil
}

func (s *specBuilderService) RegisterSecret(_ context.Context, secret *pb.Secret) (*emptypb.Empty, error) {
	if err := secret.Validate(); err != nil {
		return nil, err
	}
	s.spec.Secrets[secret.Name] = secret.Value
	return empty(), nil
}

func (s *specBuilderService) HasFunctions(_ context.Context, _ *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	var hasFunc bool

	for _, f := range s.spec.Functions {
		if hasFunc = (f.Script == ""); hasFunc {
			break
		}
	}

	return wrapperspb.Bool(hasFunc), nil
}

func (s *specBuilderService) ListResources(_ context.Context, _ *emptypb.Empty) (*pb.ListResourcesResponse, error) {
	return &pb.ListResourcesResponse{Resources: s.resources}, nil
}

func (s *specBuilderService) GetSpec(_ context.Context, req *pb.GetSpecRequest) (*pb.GetSpecResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := s.spec.SetImageForFunctions(req.Image); err != nil {
		return nil, err
	}

	if _, err := s.spec.BuildDAG(); err != nil {
		return nil, err
	}

	spec, err := s.spec.Marshal()
	if err != nil {
		return nil, err
	}

	return &pb.GetSpecResponse{Spec: spec}, nil
}

func configMap(configs *pb.Configs) map[string]any {
	if configs == nil {
		return nil
	}

	m := make(map[string]any)
	for _, c := range configs.Config {
		m[c.Field] = c.Value
	}
	return m
}
