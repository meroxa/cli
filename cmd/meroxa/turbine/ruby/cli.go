package turbinerb

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbineRbCLI struct {
	logger  log.Logger
	appPath string
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineRbCLI{logger: l, appPath: appPath}
}
