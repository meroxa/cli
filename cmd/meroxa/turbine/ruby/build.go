package turbinerb

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
)

func (t *turbineRbCLI) Build(ctx context.Context, appName string, platform bool) (string, error) {
	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandBuild,
		map[string]string{})
	err := turbine.RunCMD(ctx, t.logger, cmd)
	return t.appPath, err
}

func (t *turbineRbCLI) CleanUpBinaries(ctx context.Context, appName string) {
}
