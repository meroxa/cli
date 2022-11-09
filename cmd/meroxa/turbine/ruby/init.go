package turbinerb

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineRbCLI) Init(ctx context.Context, name string) error {
	cmd := utils.RunTurbineJS(ctx, "generate", name, t.appPath)
	_, err := utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	return err
}

func (t *turbineRbCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
