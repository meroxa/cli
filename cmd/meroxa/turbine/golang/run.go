package turbinego

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/meroxa/cli/log"
)

// Run will build a go binary and will run it.
func Run(ctx context.Context, appPath string, l log.Logger) error {
	// grab current location to use it as project name
	appName := path.Base(appPath)

	// building is a requirement prior to running for go apps
	err := BuildBinary(ctx, l, appPath, appName, false)
	if err != nil {
		return err
	}

	err = os.Chdir(appPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("./" + appName) //nolint:gosec
	output, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, string(output))
		return fmt.Errorf("run failed")
	}
	l.Info(ctx, string(output))
	return nil
}

// RunCleanup removes any dangling binaries.
func RunCleanup(path, appName string) {
	localBinary := filepath.Join(path, appName)
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
