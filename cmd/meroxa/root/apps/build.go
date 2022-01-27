package apps

import (
	"context"
	"github.com/meroxa/cli/log"
	"os/exec"
)

func buildGoApp(ctx context.Context, l log.Logger, path string, platform bool) error {
	var cmd *exec.Cmd
	if platform {
		cmd = exec.Command("go", "build", "--tags platform", path+"/...")
	} else {
		cmd = exec.Command("go", "build", path+"/...")
	}

	l.Info(ctx, "building apps...\n")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// TODO: parse output for build errors
	l.Info(ctx, "build complete!")
	l.Info(ctx, string(stdout))
	return nil
}
