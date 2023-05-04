package turbinepy

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/internal"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/client"
	"github.com/meroxa/turbine-core/pkg/server"
)

type turbinePyCLI struct {
	logger            log.Logger
	appPath           string
	grpcListenAddress string
	runner            server.Server
	builder           server.Server
	bc                client.Client
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbinePyCLI{
		logger:            l,
		appPath:           appPath,
		grpcListenAddress: internal.RandomLocalAddr(),
		runner:            server.NewRunServer(),
		builder:           server.NewSpecBuilderServer(),
	}
}
