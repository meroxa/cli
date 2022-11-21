package turbinerb

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
)

func (t *turbineRbCLI) Run(ctx context.Context) error {
	go t.runServer.Run(ctx)
	defer t.runServer.GracefulStop()

	cmd := internal.NewTurbineCmd(t.appPath, map[string]string{
		"TURBINE_CORE_SERVER": t.grpcListenAddress,
	})

	return turbine.RunCMD(ctx, t.logger, cmd)
}
