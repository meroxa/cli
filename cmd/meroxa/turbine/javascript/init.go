package turbinejs

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineJsCLI) Init(ctx context.Context, name string) error {
	cmd := RunTurbineJS(ctx, "generate", name, t.appPath)
	_, err := utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	return err
}

func (t *turbineJsCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
