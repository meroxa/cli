package turbinejs

import (
	"context"
	"fmt"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// GetVersion calls turbine-js-cli to get the version of the app.
func (t *turbineJsCLI) GetVersion(ctx context.Context) (string, error) {
	cmd := utils.RunTurbineJS(ctx, "version", t.appPath)
	stdOut, err := utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)

	version, err := utils.GetTurbineResponseFromOutput(stdOut)
	if err != nil {
		err = fmt.Errorf(
			"unable to determine the version of turbine-js-cli used for the Meroxa Application at %s", t.appPath)
	}
	return version, err
}
