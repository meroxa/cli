package turbinego

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbineGoCLI struct {
	logger  log.Logger
	appPath string
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineGoCLI{logger: l, appPath: appPath}
}
