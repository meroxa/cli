package turbinejs

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/javascript/internal"
)

// Run calls turbine-js to exercise the app.
func (t *turbineJsCLI) Run(ctx context.Context) error {
	go t.runner.RunAddr(ctx, t.grpcListenAddress)
	defer t.runner.GracefulStop()

	cmd := internal.NewTurbineCmd(
		ctx,
		t.appPath,
		internal.TurbineCommandRun,
		map[string]string{
			"TURBINE_CORE_SERVER": t.grpcListenAddress,
		},
		"devel",
	)

	return utils.RunCMD(ctx, t.logger, cmd)
}
