//go:generate mockgen -source=cli.go -package=mock -destination=mock/cli_mock.go

package turbinerb

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/server"
	"google.golang.org/grpc"
)

type turbineServer interface {
	Run(context.Context)
	GracefulStop()
}

type recordClient struct {
	*grpc.ClientConn
	pb.TurbineServiceClient
}

type turbineRbCLI struct {
	logger            log.Logger
	appPath           string
	runServer         turbineServer
	recordServer      turbineServer
	recordClient      recordClient
	grpcListenAddress string
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineRbCLI{
		logger:            l,
		appPath:           appPath,
		runServer:         server.NewRunServer(),
		recordServer:      server.NewRecordServer(),
		grpcListenAddress: server.ListenAddress,
	}
}
