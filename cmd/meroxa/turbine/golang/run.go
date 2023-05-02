package turbinego

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

// Build will create a go binary with a specific name on a specific path.
func (t *turbineGoCLI) build(ctx context.Context, appName string) error {
	var cmd *exec.Cmd

	t.logger.StartSpinner("\t", "Building Golang binary...")
	cmd = exec.CommandContext(ctx, "go", "build", "-o", appName, "./...")
	cmd.Dir = t.appPath

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		t.logger.StopSpinnerWithStatus(string(stdout), log.Failed)
		return fmt.Errorf("build failed")
	}
	return t.buildForCrossCompile(ctx, appName)
}

func (t *turbineGoCLI) buildForCrossCompile(ctx context.Context, appName string) error {
	cmd := exec.CommandContext(ctx, "go", "build", "-o", appName+".cross", "./...") //nolint:gosec
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	cmd.Dir = t.appPath

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		t.logger.StopSpinnerWithStatus(string(stdout), log.Failed)
		return fmt.Errorf("cross compile failed")
	}
	t.logger.StopSpinnerWithStatus("Successfully built the Golang binary", log.Successful)
	return nil
}

// Run will build a go binary and will run it.
func (t *turbineGoCLI) Run(ctx context.Context) error {
	appName, err := utils.GetAppNameFromAppJSON(t.logger, t.appPath)
	if err != nil {
		return err
	}

	// building is a requirement prior to running for go apps
	if err = t.build(ctx, appName); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, t.appPath+"/"+appName) //nolint:gosec
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.logger.Error(ctx, string(output))
		return fmt.Errorf("run failed")
	}
	t.logger.Info(ctx, string(output))
	return nil
}

// RunCleanup removes any dangling binaries.
// TODO: Remove once `apps run` for go uses turbine-core exclusively.
func RunCleanup(ctx context.Context, l log.Logger, appPath, appName string) {
	localBinary := filepath.Join(appPath, appName)
	err := os.Remove(localBinary)
	if err != nil {
		l.Warnf(ctx, "warning: failed to clean up %s\n", localBinary)
	}

	crossCompiledBinary := localBinary + ".cross"

	err = os.Remove(crossCompiledBinary)
	if err != nil {
		l.Warnf(ctx, "warning: failed to clean up %s\n", crossCompiledBinary)
	}
}
