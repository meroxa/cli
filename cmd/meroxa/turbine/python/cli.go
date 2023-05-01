package turbinepy

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/server"
)

type turbinePyCLI struct {
	logger            log.Logger
	appPath           string
	tmpPath           string
	grpcListenAddress string
	runner            server.Server
	builder           server.Server
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbinePyCLI{
		logger:            l,
		appPath:           appPath,
		grpcListenAddress: server.ListenAddress,
		runner:            server.NewRunServer(),
		builder:           server.NewSpecBuilderServer(),
	}
}
