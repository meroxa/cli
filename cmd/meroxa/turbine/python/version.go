package turbinepy

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// GetVersion calls turbine-py to get the version of the app.
func (t *turbinePyCLI) GetVersion(ctx context.Context) (string, error) {
	cmd := exec.Command("turbine-py", "version", t.appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(
			"unable to determine the turbine-py version of the Meroxa Application at %s",
			t.appPath)
		return "", err
	}

	version, err := utils.GetTurbineResponseFromOutput(string(output))
	if err != nil {
		// For versions <= v1.6.1
		version = string(output)
		version = strings.TrimSpace(version)
		return version, nil
	}
	return version, err
}
