package apps

import (
	"context"
	"os/exec"

	"github.com/meroxa/cli/log"
)

func buildGoApp(ctx context.Context, l log.Logger, appPath, appName string, platform bool) error {
	var cmd *exec.Cmd
	if platform {
		cmd = exec.Command("go", "build", "--tags", "platform", "-o", "./"+appPath, "./"+appPath+"/...") //nolint:gosec
	} else {
		cmd = exec.Command("go", "build", "-o", appName, "./...")
	}
	cmd.Dir = appPath

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, string(stdout))
		return err
	}

	return nil
}
