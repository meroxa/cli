package turbinego

import (
	"context"
	"fmt"
	"os"
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
	return buildForCrossCompile(ctx, l, appPath, appName)
}

func buildForCrossCompile(ctx context.Context, l log.Logger, appPath, appName string) error {
	cmd := exec.Command("go", "build", "--tags", "platform", "-o", appName+".cross", "./...") //nolint:gosec
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	cmd.Dir = appPath

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Errorf(ctx, "%s", string(stdout))
		return fmt.Errorf("cross compile failed")
	}
	return nil
}
