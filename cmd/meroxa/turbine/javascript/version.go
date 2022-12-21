package turbinejs

import (
	"context"
	"fmt"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// GetVersion calls turbine-js-cli to get the turbine-js-framework version used by the given app.
func (t *turbineJsCLI) GetVersion(ctx context.Context) (string, error) {
	stdOut, err := RunTurbineJS(ctx, t.logger, "version", t.appPath)
	fmtErr := fmt.Errorf(
		"unable to determine the version of turbine-js-framework used by the Meroxa Application at %s: %s",
		t.appPath, err)
	if err != nil {
		return "", fmtErr
	}

	version, err := utils.GetTurbineResponseFromOutput(stdOut)
	if err != nil {
		err = fmtErr
	}
	return version, err
}
