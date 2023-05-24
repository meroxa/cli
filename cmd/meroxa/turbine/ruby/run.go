package turbinerb

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
)

func (t *turbineRbCLI) Run(ctx context.Context) error {
	gracefulStop, grpcListenAddress := t.Core.Run(ctx)
	defer gracefulStop()

	cmd := internal.NewTurbineCmd(
		ctx,
		t.appPath,
		internal.TurbineCommandRun,
		map[string]string{
			"TURBINE_CORE_SERVER": grpcListenAddress,
		})
	return turbine.RunCMD(ctx, t.logger, cmd)
}
