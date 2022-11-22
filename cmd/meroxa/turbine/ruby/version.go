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
		"-e", `'puts Gem.loaded_specs["turbin_rb"].version.version'`,
	}...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"unable to determine the turbine-rb version of the Meroxa Application at %s",
			t.appPath)
	}

	return strings.TrimSpace(string(output)), nil
}
