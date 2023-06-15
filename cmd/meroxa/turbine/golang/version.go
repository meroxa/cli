package turbinego

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// GetVersion will return the tag or hash of the turbine-go dependency of a given app.
func (t *turbineGoCLI) GetVersion(ctx context.Context) (string, error) {
	var cmd *exec.Cmd

	cmd = exec.CommandContext(
		ctx,
		"go",
		"list", "-m", "-f", "'{{ .Version }}'", "github.com/meroxa/turbine-go/v2")
	cmd.Dir = t.appPath
	fmtErr := fmt.Errorf(
		"unable to determine the version of turbine-go used by the Meroxa Application at %s",
		t.appPath)

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmtErr
	}

	version := strings.TrimSpace(string(stdout))
	chars := []rune(version)
	if chars[0] == 'v' {
		// Looks like v0.0.0-20221024132549-e6470e58b719
		const sections = 3
		parts := strings.Split(version, "-")
		if len(parts) < sections {
			return "", fmtErr
		}
		version = parts[2]
	}
	return version, nil
}
