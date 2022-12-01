package turbinepy

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbinePyCLI struct {
	logger  log.Logger
	appPath string
	tmpPath string
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbinePyCLI{logger: l, appPath: appPath}
}
