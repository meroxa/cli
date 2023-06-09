package turbinepy

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/python/internal"
)

func (t *turbinePyCLI) Run(ctx context.Context) error {
	gracefulStop, grpcListenAddress := t.Core.Run(ctx)
	defer gracefulStop()

	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandRun,
		map[string]string{
			"TURBINE_CORE_SERVER": grpcListenAddress,
		},
		t.appPath,
		"",
	)
	_, err := turbine.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	return err
}
