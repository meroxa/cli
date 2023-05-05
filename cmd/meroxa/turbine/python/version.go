package turbinepy

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// GetVersion calls turbine-py to get its version used by the given app.
func (t *turbinePyCLI) GetVersion(_ context.Context) (string, error) {
	cmd := exec.Command("turbine-py", "version", t.appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(
			"unable to determine the version of turbine-py used by the Meroxa Application at %s",
			t.appPath)
		return "", err
	}

	version, err := utils.GetTurbineResponseFromOutput(string(output))
	// FYI the error return here is being used to fallback to previous CLI output format
	if err != nil {
		// For versions <= v1.6.1
		version = string(output)
		version = strings.TrimSpace(version)
		return version, nil
	}

	return version, err
}
