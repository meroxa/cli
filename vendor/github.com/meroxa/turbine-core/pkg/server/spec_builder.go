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

	spec *ir.DeploymentSpec
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

func (s *specBuilderService) ReadFromSource(_ context.Context, req *pb.ReadFromSourceRequest) (*pb.RecordsCollection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	c := ir.ConnectorSpec{
		UUID:       uuid.New().String(),
		PluginName: req.PluginName,
		Direction:  ir.ConnectorSource,
		Config:     configMap(req.Configuration),
	}

	if err := s.spec.AddSource(&c); err != nil {
		return nil, err
	}

	return &pb.RecordsCollection{
		Stream: c.UUID,
	}, nil
}

func (s *specBuilderService) WriteToDestination(_ context.Context, req *pb.WriteToDestinationRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	c := ir.ConnectorSpec{
		UUID:       uuid.New().String(),
		PluginName: req.PluginName,
		Direction:  ir.ConnectorDestination,
		Config:     configMap(req.Configuration),
	}
	if err := s.spec.AddDestination(&c); err != nil {
		return nil, err
	}

	if err := s.spec.AddStream(&ir.StreamSpec{
		UUID:     uuid.New().String(),
		FromUUID: req.Records.Stream,
		ToUUID:   c.UUID,
		Name:     req.Records.Stream + "_" + c.UUID,
	}); err != nil {
		return nil, err
	}

	return empty(), nil
}

func (s *specBuilderService) Process(_ context.Context, req *pb.ProcessRecordsRequest) (*pb.RecordsCollection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	f := ir.FunctionSpec{
		UUID: uuid.New().String(),
		Name: strings.ToLower(req.Process.Name),
	}
	if err := s.spec.AddFunction(&f); err != nil {
		return nil, err
	}

	if err := s.spec.AddStream(&ir.StreamSpec{
		UUID:     uuid.New().String(),
		FromUUID: req.Records.Stream,
		ToUUID:   f.UUID,
		Name:     req.Records.Stream + "_" + f.UUID,
	}); err != nil {
		return nil, err
	}

	return &pb.RecordsCollection{
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
	return wrapperspb.Bool(len(s.spec.Functions) > 0), nil
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

func configMap(configs *pb.Configurations) map[string]any {
	if configs == nil {
		return nil
	}

	m := make(map[string]any)
	for _, c := range configs.Configuration {
		m[c.Field] = c.Value
	}
	return m
}
