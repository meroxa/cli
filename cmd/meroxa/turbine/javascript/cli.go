package turbinejs

import (
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbineJsCLI struct {
	logger log.Logger
}

func New(l log.Logger) turbine.CLI {
	return &turbineJsCLI{logger: l}
}
