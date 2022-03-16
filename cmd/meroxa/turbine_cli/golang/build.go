package turbinego

import (
	"context"
	"os/exec"

	"github.com/meroxa/cli/log"
)

// buildBinary will create a go binary with a specific name on a specific path.
func buildBinary(ctx context.Context, l log.Logger, appPath, appName string, platform bool) error {
	var cmd *exec.Cmd

	if platform {
		cmd = exec.Command("go", "build", "--tags", "platform", "-o", appName, "./...")
	} else {
		cmd = exec.Command("go", "build", "-o", appName, "./...")
	}
	cmd.Dir = appPath

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Errorf(ctx, string(stdout))
		return err
	}

	return nil
}
