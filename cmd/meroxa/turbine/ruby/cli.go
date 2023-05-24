//go:generate mockgen -source=cli.go -package=mock -destination=mock/cli_mock.go

package turbinerb

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbineRbCLI struct {
	logger  log.Logger
	appPath string
	*turbine.Core
	*turbine.Git
	*turbine.Docker
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineRbCLI{
		logger:  l,
		appPath: appPath,
		Core:    turbine.NewCore(),
	}
}
