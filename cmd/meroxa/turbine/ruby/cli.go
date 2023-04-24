//go:generate mockgen -source=cli.go -package=mock -destination=mock/cli_mock.go

package turbinerb

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"

	"github.com/meroxa/turbine-core/pkg/client"
	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/meroxa/turbine-core/pkg/server"
)

type turbineRbCLI struct {
	logger            log.Logger
	appPath           string
	grpcListenAddress string

	builder server.Server
	runner  server.Server
	bc      client.Client
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineRbCLI{
		logger:            l,
		appPath:           appPath,
		runner:            server.NewRunServer(),
		builder:           server.NewSpecBuilderServer(),
		grpcListenAddress: server.ListenAddress,
	}
}
