package turbinepy

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/python/internal"
)

func (t *turbinePyCLI) CreateDockerfile(ctx context.Context, _ string) (string, error) {
	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandBuild,
		map[string]string{},
		t.appPath,
	)
	_, err := turbine.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	return t.appPath, err
}

func (t *turbinePyCLI) StartGrpcServer(ctx context.Context, gitSha string) (func(), error) {
	grpcListenAddress, err := t.Core.Start(ctx)
	if err != nil {
		return nil, err
	}

	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandRecord,
		map[string]string{
			"TURBINE_CORE_SERVER": grpcListenAddress,
		},
		t.appPath,
		gitSha,
	)

	if _, err := turbine.RunCmdWithErrorDetection(ctx, cmd, t.logger); err != nil {
		return nil, err
	}

	return t.Core.Stop()
}
