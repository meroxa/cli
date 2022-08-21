package turbinego

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbineGoCLI struct {
	logger log.Logger
}

func New(l log.Logger) turbine.CLI {
	return &turbineGoCLI{logger: l}
}
