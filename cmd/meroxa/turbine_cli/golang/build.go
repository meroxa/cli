package turbinego

import (
	"context"
	"os/exec"

	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"

	"github.com/meroxa/cli/log"
)

// BuildBinary will create a go binary with a specific name on a specific path.
func BuildBinary(ctx context.Context, l log.Logger, appPath, appName string, platform bool) error {
	var cmd *exec.Cmd

	if platform {
		cmd = exec.Command("go", "build", "--tags", "platform", "-o", appName, "./...")
	} else {
		cmd = exec.Command("go", "build", "-o", appName, "./...")
	}
	cmd.Dir = appPath

	_, err := turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}
