package turbinepy

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/python/internal"
)

func (t *turbinePyCLI) Run(ctx context.Context) error {
	go t.runner.Run(ctx)
	defer t.runner.GracefulStop()
	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandRun,
		map[string]string{
			"TURBINE_CORE_SERVER": t.grpcListenAddress,
		},
		t.appPath)
	return turbine.RunCMD(ctx, t.logger, cmd)
}
