package server

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/meroxa/turbine-core/v2/pkg/ir"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ turbinev2.ServiceServer = (*SpecBuilderService)(nil)

type SpecBuilderService struct {
	turbinev2.UnimplementedServiceServer

	spec *ir.DeploymentSpec
	// resources []*turbinev2.Resource
}

func NewSpecBuilderService() *SpecBuilderService {
	return &SpecBuilderService{
		spec: &ir.DeploymentSpec{},
	}
}

func (s *SpecBuilderService) Init(_ context.Context, req *turbinev2.InitRequest) (*emptypb.Empty, error) {
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

func (s *SpecBuilderService) AddSource(_ context.Context, req *turbinev2.AddSourceRequest) (*turbinev2.AddSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	c := ir.ConnectorSpec{
		UUID:         uuid.New().String(),
		Name:         req.Name,
		PluginType:   ir.PluginSource,
		PluginName:   req.Plugin.Name,
		PluginConfig: req.Plugin.Config,
	}

	if err := s.spec.AddSource(&c); err != nil {
		return nil, err
	}

	return &turbinev2.AddSourceResponse{StreamName: c.UUID}, nil
}

func (s *SpecBuilderService) ReadRecords(_ context.Context, req *turbinev2.ReadRecordsRequest) (*turbinev2.ReadRecordsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &turbinev2.ReadRecordsResponse{
		StreamRecords: &turbinev2.StreamRecords{
			StreamName: req.SourceStream,
		},
	}, nil
}

func (s *SpecBuilderService) AddDestination(_ context.Context, req *turbinev2.AddDestinationRequest) (*turbinev2.AddDestinationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	c := ir.ConnectorSpec{
		UUID:         uuid.New().String(),
		Name:         req.Name,
		PluginType:   ir.PluginDestination,
		PluginName:   req.Plugin.Name,
		PluginConfig: req.Plugin.Config,
	}

	if err := s.spec.AddDestination(&c); err != nil {
		return nil, err
	}

	return &turbinev2.AddDestinationResponse{
		Id: c.UUID,
	}, nil
}

func (s *SpecBuilderService) WriteRecords(_ context.Context, req *turbinev2.WriteRecordsRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := s.spec.AddStream(&ir.StreamSpec{
		UUID:     uuid.New().String(),
		FromUUID: req.StreamRecords.StreamName,
		ToUUID:   req.DestinationID,
		Name:     req.StreamRecords.StreamName + "_" + req.DestinationID,
	}); err != nil {
		return nil, err
	}

	return empty(), nil
}

func (s *SpecBuilderService) ProcessRecords(_ context.Context, req *turbinev2.ProcessRecordsRequest) (*turbinev2.ProcessRecordsResponse, error) {
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
		FromUUID: req.StreamRecords.StreamName,
		ToUUID:   f.UUID,
		Name:     req.StreamRecords.StreamName + "_" + f.UUID,
	}); err != nil {
		return nil, err
	}

	return &turbinev2.ProcessRecordsResponse{
		StreamRecords: &turbinev2.StreamRecords{
			StreamName: f.UUID,
			Records:    req.StreamRecords.Records,
		},
	}, nil
}

func (s *SpecBuilderService) GetSpec(_ context.Context, req *turbinev2.GetSpecRequest) (*turbinev2.GetSpecResponse, error) {
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

	return &turbinev2.GetSpecResponse{Spec: spec}, nil
}
