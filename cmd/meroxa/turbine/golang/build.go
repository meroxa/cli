package turbinego

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/meroxa/cli/log"
)

// Build will create a go binary with a specific name on a specific path.
func (t *turbineGoCLI) Build(ctx context.Context, appName string, platform bool) (string, error) {
	var cmd *exec.Cmd

	t.logger.StartSpinner("\t", " Building Golang binary...")
	if platform {
		cmd = exec.CommandContext(ctx, "go", "build", "--tags", "platform", "-o", appName, "./...")
	} else {
		cmd = exec.CommandContext(ctx, "go", "build", "-o", appName, "./...")
	}
	cmd.Dir = t.appPath

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		t.logger.StopSpinnerWithStatus(string(stdout), log.Failed)
		return "", fmt.Errorf("build failed")
	}
	return "", t.buildForCrossCompile(ctx, appName)
}

func (t *turbineGoCLI) buildForCrossCompile(ctx context.Context, appName string) error {
	cmd := exec.CommandContext(ctx, "go", "build", "--tags", "platform", "-o", appName+".cross", "./...") //nolint:gosec
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	cmd.Dir = t.appPath

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		t.logger.StopSpinnerWithStatus(string(stdout), log.Failed)
		return fmt.Errorf("cross compile failed")
	}
	t.logger.StopSpinnerWithStatus("Successfully built the Golang binary!", log.Successful)
	return nil
}

func (t *turbineGoCLI) CleanUpBinaries(ctx context.Context, appName string) {
	RunCleanup(ctx, t.logger, t.appPath, appName)
}
