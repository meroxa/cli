package turbinego

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// Run will build a go binary and will run it.
func (t *turbineGoCLI) Run(ctx context.Context) error {
	// grab current location to use it as project name
	appName, err := utils.GetAppNameFromAppJSON(ctx, t.logger, t.appPath)
	if err != nil {
		return err
	}

	// building is a requirement prior to running for go apps
	err = BuildBinary(ctx, t.logger, t.appPath, appName, false)
	if err != nil {
		return err
	}

	cmd := exec.Command("./" + appName) //nolint:gosec
	cmd.Path = t.appPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.logger.Error(ctx, string(output))
		return fmt.Errorf("run failed")
	}
	t.logger.Info(ctx, string(output))
	return nil
}

// RunCleanup removes any dangling binaries.
func RunCleanup(appPath, appName string) {
	localBinary := filepath.Join(appPath, appName)
	err := os.Remove(localBinary)
	if err != nil {
		fmt.Printf("warning: failed to clean up %s\n", localBinary)
	}

	crossCompiledBinary := localBinary + ".cross"
	err = os.Remove(crossCompiledBinary)
	if err != nil {
		fmt.Printf("warning: failed to clean up %s\n", crossCompiledBinary)
	}
}
