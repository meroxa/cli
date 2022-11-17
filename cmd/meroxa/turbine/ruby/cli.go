//go:generate mockgen -source=cli.go -package=mock -destination=mock/cli_mock.go

package turbinerb

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/server"
)

type turbineServer interface {
	Run(context.Context)
	GracefulStop()
}

type turbineRbCLI struct {
	logger            log.Logger
	appPath           string
	runServer         turbineServer
	grpcListenAddress string
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineRbCLI{
		logger:            l,
		appPath:           appPath,
		runServer:         server.NewRunServer(),
		grpcListenAddress: server.ListenAddress,
	}
}
