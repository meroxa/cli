package turbinejs

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/javascript/internal"
)

// Run calls turbine-js to exercise the app.
func (t *turbineJsCLI) Run(ctx context.Context) error {
	gracefulStop, grpcListenAddress := t.Core.Run(ctx)
	defer gracefulStop()

	cmd := internal.NewTurbineCmd(
		ctx,
		t.appPath,
		internal.TurbineCommandRun,
		map[string]string{
			"TURBINE_CORE_SERVER": grpcListenAddress,
		},
		"devel",
	)

	_, err := utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	return err
}
