package turbinejava

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/turbine/java/internal"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineJavaCLI) CreateDockerfile(ctx context.Context, _ string) (string, error) {
	cmd := internal.NewTurbineCmd(
		ctx,
		t.appPath,
		internal.TurbineCommandBuild,
		map[string]string{},
		t.appPath,
	)
	return t.appPath, turbine.RunCMD(ctx, t.logger, cmd)
}

func (t *turbineJavaCLI) StartGrpcServer(ctx context.Context, gitSha string) (func(), error) {
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

	if err := turbine.RunCMD(ctx, t.logger, cmd); err != nil {
		return nil, err
	}

	return t.Core.Stop()
}
