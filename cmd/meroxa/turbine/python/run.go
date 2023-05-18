package turbinepy

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/python/internal"
)

func (t *turbinePyCLI) Run(ctx context.Context) error {
	go t.runner.RunAddr(ctx, t.grpcListenAddress)
	defer t.runner.GracefulStop()
	gitSha, err := t.GetGitSha(ctx)
	if err != nil {
		return err
	}
	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandRun,
		map[string]string{
			"TURBINE_CORE_SERVER": t.grpcListenAddress,
		},
		t.appPath,
		gitSha,
	)
	return turbine.RunCMD(ctx, t.logger, cmd)
}
