package turbinego

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/internal"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/client"
	"github.com/meroxa/turbine-core/pkg/server"
)

type turbineGoCLI struct {
	logger            log.Logger
	appPath           string
	grpcListenAddress string

	builder server.Server
	runner  server.Server
	bc      client.Client
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineGoCLI{
		logger:            l,
		appPath:           appPath,
		runner:            server.NewRunServer(),
		builder:           server.NewSpecBuilderServer(),
		grpcListenAddress: internal.RandomLocalAddr(),
	}
}
