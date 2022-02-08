package apps

import (
	"context"
	"os/exec"

	"github.com/meroxa/cli/log"
)

// This will be used for running apps locally.
func buildGoApp(ctx context.Context, l log.Logger, path string, platform bool) error {
	var cmd *exec.Cmd
	if platform {
		cmd = exec.Command("go", "build", "--tags", "platform", "-o", "./"+path, "./"+path+"/...") //nolint:gosec
	} else {
		cmd = exec.Command("go", "build", path+"/...") //nolint:gosec
	}

	l.Infof(ctx, "Building app located at %q ...\n", path)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, string(stdout))
		return err
	}
	// TODO: parse output for build errors
	l.Info(ctx, "Build complete!")
	return nil
}
