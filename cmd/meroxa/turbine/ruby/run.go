package turbinerb

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
)

func (t *turbineRbCLI) Run(ctx context.Context) error {
	go t.runner.Run(ctx)
	defer t.runner.GracefulStop()

	cmd := internal.NewTurbineCmd(
		ctx,
		t.appPath,
		internal.TurbineCommandRun,
		map[string]string{
			"TURBINE_CORE_SERVER": t.grpcListenAddress,
		})
	return turbine.RunCMD(ctx, t.logger, cmd)
}
