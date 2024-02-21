package server

import (
	"context"
	"fmt"
	"path"

	"github.com/meroxa/turbine-core/v2/pkg/app"
	"github.com/meroxa/turbine-core/v2/pkg/server/internal"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ turbinev2.ServiceServer = (*RunService)(nil)

type RunService struct {
	turbinev2.UnimplementedServiceServer

	config  app.Config
	appPath string
}

func NewRunService() *RunService {
	return &RunService{
		config: app.Config{},
	}
}

func (s *RunService) Init(_ context.Context, req *turbinev2.InitRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	config, err := app.ReadConfig(req.AppName, req.ConfigFilePath)
	if err != nil {
		return nil, err
	}
	s.config = config
	s.appPath = req.ConfigFilePath

	return empty(), nil
}

func (s *RunService) AddSource(_ context.Context, req *turbinev2.AddSourceRequest) (*turbinev2.AddSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &turbinev2.AddSourceResponse{
		StreamName: req.Name,
	}, nil
}

func (s *RunService) ReadRecords(ctx context.Context, req *turbinev2.ReadRecordsRequest) (*turbinev2.ReadRecordsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	fixtureFile, ok := s.config.Fixtures[req.SourceStream]
	if !ok {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf(
				"no fixture file found for source %s. Ensure that the source is declared in your app.json.",
				req.SourceStream,
			),
		)
	}

	rr, err := internal.ReadFixture(ctx, path.Join(s.appPath, fixtureFile))
	if err != nil {
		return nil, err
	}

	return &turbinev2.ReadRecordsResponse{
		StreamRecords: &turbinev2.StreamRecords{
			StreamName: req.SourceStream,
			Records:    rr,
		},
	}, nil
}

func (s *RunService) AddDestination(_ context.Context, req *turbinev2.AddDestinationRequest) (*turbinev2.AddDestinationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &turbinev2.AddDestinationResponse{
		Id: req.Name,
	}, nil
}

func (s *RunService) WriteRecords(_ context.Context, req *turbinev2.WriteRecordsRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	internal.PrintRecords(req.DestinationID, req.StreamRecords)

	return empty(), nil
}

func (s *RunService) ProcessRecords(_ context.Context, req *turbinev2.ProcessRecordsRequest) (*turbinev2.ProcessRecordsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &turbinev2.ProcessRecordsResponse{
		StreamRecords: &turbinev2.StreamRecords{
			StreamName: req.StreamRecords.StreamName,
			Records:    req.StreamRecords.Records,
		},
	}, nil
}
