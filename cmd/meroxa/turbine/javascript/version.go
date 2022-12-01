package turbinejs

import (
	"context"
	"fmt"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// GetVersion calls turbine-js-cli to get the turbine-js-framework version used by the given app.
func (t *turbineJsCLI) GetVersion(ctx context.Context) (string, error) {
	cmd := RunTurbineJS(ctx, "version", t.appPath)
	stdOut, err := utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	fmtErr := fmt.Errorf(
		"unable to determine the version of turbine-js-framework used by the Meroxa Application at %s",
		t.appPath)
	if err != nil {
		return "", fmtErr
	}

	version, err := utils.GetTurbineResponseFromOutput(stdOut)
	if err != nil {
		err = fmtErr
	}
	return version, err
}
