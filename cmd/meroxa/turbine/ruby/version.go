package turbinerb

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// GetVersion will return the version of turbine-rb dependency of a given app.
func (t *turbineRbCLI) GetVersion(ctx context.Context) (string, error) {
	cmd := exec.Command("ruby", []string{
		"-r", "rubygems",
		"-r", fmt.Sprintf("%s/app", t.appPath),
		"-e", `puts Gem.loaded_specs["turbine"].version.version`,
	}...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(
			"unable to determine the turbine-rb version of the Meroxa Application at %s",
			t.appPath)
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
