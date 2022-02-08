package apps

import (
	"context"
	"os/exec"

	"github.com/meroxa/cli/log"
)

func buildGoApp(ctx context.Context, l log.Logger, path string, platform bool) error {
	var cmd *exec.Cmd
	if platform {
		cmd = exec.Command("go", "build", "--tags", "platform", "-o", "./"+path, "./"+path+"/...") //nolint:gosec
	} else {
		cmd = exec.Command("go", "build", path+"/...") //nolint:gosec
	}

	l.Info(ctx, "building app...\n")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, string(stdout))
		return err
	}
	// TODO: parse output for build errors
	l.Info(ctx, "build complete!")
	return nil
}
