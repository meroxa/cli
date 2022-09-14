package turbinejs

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// Run calls turbine-js to exercise the app.
func (t *turbineJsCLI) Run(ctx context.Context) (err error) {
	cmd := utils.RunTurbineJS(ctx, "test", t.appPath)
	stdOut, err := utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	t.logger.Info(ctx, stdOut)
	return err
}
