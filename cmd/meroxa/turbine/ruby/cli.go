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

type specBuilderClient struct {
	*grpc.ClientConn
	pb.TurbineServiceClient
}

type turbineRbCLI struct {
	logger            log.Logger
	appPath           string
	grpcListenAddress string

	bs turbineServer
	bc specBuilderClient
	rs turbineServer
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineRbCLI{
		logger:            l,
		appPath:           appPath,
		rs:                server.NewRunServer(),
		bs:                server.NewSpecBuilderServer(),
		grpcListenAddress: server.ListenAddress,
	}
}
