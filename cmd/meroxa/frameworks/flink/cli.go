package flink

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type flinkCLI struct {
	logger  log.Logger
	appPath string
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &flinkCLI{logger: l, appPath: appPath}
}
