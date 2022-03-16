package turbinejs

import (
	"context"
	"os/exec"

	"github.com/meroxa/cli/log"
)

// Build calls turbine-js to build an application.
func Build(ctx context.Context, l log.Logger) error {
	// TODO: Handle the requirement of https://github.com/meroxa/turbine-js.git being installed
	// cd into the path first
	cmd := exec.Command("npx", "turbine", "test")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	l.Info(ctx, string(stdout))
	return nil
}
