package turbinejs

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineJsCLI) Init(ctx context.Context, name string) error {
	_, err := RunTurbineJS(ctx, t.logger, "generate", name, t.appPath)
	return err
}

func (t *turbineJsCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
