package turbinejava

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbineJavaCLI struct {
	logger  log.Logger
	appPath string
	*turbine.Core
	*turbine.Git
	*turbine.Docker
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineJavaCLI{
		logger:  l,
		appPath: appPath,
		Core:    turbine.NewCore(),
	}
}
