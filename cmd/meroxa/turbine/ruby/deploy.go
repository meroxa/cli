package turbinerb

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
)

func (t *turbineRbCLI) CreateDockerfile(ctx context.Context, _ string) (string, error) {
	cmd := internal.NewTurbineCmd(
		ctx,
		t.appPath,
		internal.TurbineCommandBuild,
		map[string]string{},
	)
	_, err := turbine.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	return t.appPath, err
}

func (t *turbineRbCLI) StartGrpcServer(ctx context.Context, gitSha string) (func(), error) {
	grpcListenAddress, err := t.Core.Start(ctx)
	if err != nil {
		return nil, err
	}

	cmd := internal.NewTurbineCmd(
		ctx,
		t.appPath,
		internal.TurbineCommandRecord,
		map[string]string{
			"TURBINE_CORE_SERVER": grpcListenAddress,
		},
		gitSha,
	)

	if _, err := turbine.RunCmdWithErrorDetection(ctx, cmd, t.logger); err != nil {
		return nil, err
	}

	return t.Core.Stop()
}
