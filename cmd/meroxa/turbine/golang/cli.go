package turbinego

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

type turbineGoCLI struct {
	logger  log.Logger
	appPath string
	core    *turbine.CoreV2

	*turbine.Git
	*turbine.Docker
}

func New(l log.Logger, appPath string) turbine.CLI {
	return &turbineGoCLI{
		logger:  l,
		appPath: appPath,
		core:    turbine.NewCoreV2(),
	}
}

func (t *turbineGoCLI) GetDeploymentSpec(ctx context.Context, image string) (string, error) {
	return t.core.GetDeploymentSpec(ctx, image)
}
