package turbinego

import (
	"context"
	"fmt"
	"os/exec"

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

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Errorf(ctx, "%s", string(stdout))
		return fmt.Errorf("build failed")
	}
	return err
}
