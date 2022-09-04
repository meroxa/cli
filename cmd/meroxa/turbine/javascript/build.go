package turbinejs

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// Build calls turbine-js to build an application.
func (t *turbineJsCLI) Build(ctx context.Context, appName string, platform bool) (string, error) {
	// TODO: Handle the requirement of https://github.com/meroxa/turbine-js.git being installed
	// cd into the path first
	cmd := utils.RunTurbineJS(ctx, "test", t.appPath)
	stdOut, err := utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	t.logger.Info(ctx, stdOut)
	return "", err
}

func (t *turbineJsCLI) CleanUpTempBuildLocation(ctx context.Context) error {
	return nil
}

func (t *turbineJsCLI) CleanUpBinaries(ctx context.Context, appName string) {
}
