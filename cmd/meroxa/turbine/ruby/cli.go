//go:generate mockgen -source=cli.go -package=mock -destination=mock/cli_mock.go

package turbinerb

import (
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/server"
)

var (
	TurbineRubyFeatureFlag    = "ruby_implementation"
	ErrTurbineRubyFeatureFlag = fmt.Errorf(`no access to the Meroxa Turbine Ruby feature.
Sign up for the Beta here: https://share.hsforms.com/1Uq6UYoL8Q6eV5QzSiyIQkAc2sme`)
)

type turbineRbCLI struct {
	logger            log.Logger
	appPath           string
	runServer         turbine.TurbineServer
	recordServer      turbine.TurbineServer
	recordClient      turbine.RecordClient
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
