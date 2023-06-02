package server

import (
	"context"
	"fmt"
	"path"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/server/internal"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pb.TurbineServiceServer = (*runService)(nil)

type runService struct {
	pb.UnimplementedTurbineServiceServer

	config  app.Config
	appPath string
}

func NewRunService() *runService {
	return &runService{
		config: app.Config{},
	}
}

func (s *runService) Init(ctx context.Context, req *pb.InitRequest) (*emptypb.Empty, error) {
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

func (s *runService) ReadFromSource(ctx context.Context, req *pb.ReadFromSourceRequest) (*pb.RecordsCollection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	fixtureFile, ok := s.config.Fixtures[req.PluginName]
	if !ok {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf(
				"No fixture file found for plugin name %s. Ensure that the plugin name is declared in your app.json.",
				req.PluginName,
			),
		)
	}

	fixture := &internal.FixtureConnector{
		PluginName: req.PluginName,
		File: path.Join(
			s.appPath,
			fixtureFile,
		),
	}

	rr, err := fixture.ReadAll(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.RecordsCollection{
		Records: rr,
	}, nil
}

func (s *runService) WriteToDestination(_ context.Context, req *pb.WriteToDestinationRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	internal.PrintRecords(req.PluginName, req.Records.Records)

	return empty(), nil
}

func (s *runService) Process(ctx context.Context, req *pb.ProcessRecordsRequest) (*pb.RecordsCollection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return req.Records, nil
}

func (s *runService) RegisterSecret(ctx context.Context, req *pb.Secret) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return empty(), nil
}
