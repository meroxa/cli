package turbinejs

import (
	"context"

	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"

	"github.com/meroxa/cli/log"
)

// Build calls turbine-js to build an application.
func Build(ctx context.Context, l log.Logger, path string) error {
	// TODO: Handle the requirement of https://github.com/meroxa/turbine-js.git being installed
	// cd into the path first
	cmd := turbinecli.RunTurbineJS("test", path)
	stdOut, err := turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
	l.Info(ctx, stdOut)
	return err
}
