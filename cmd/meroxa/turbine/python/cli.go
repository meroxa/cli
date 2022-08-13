package turbinepy

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbinePyCLI struct {
	logger log.Logger
}

func New(l log.Logger) turbine.CLI {
	return &turbinePyCLI{logger: l}
}
