package turbinejs

import (
	"context"
	"os/exec"

	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"

	"github.com/meroxa/cli/log"
)

// Build calls turbine-js to build an application.
func Build(ctx context.Context, l log.Logger, path string) error {
	// TODO: Handle the requirement of https://github.com/meroxa/turbine-js.git being installed
	// cd into the path first
	cmd := exec.Command("npx", "turbine", "test", path)
	_, err := turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}
