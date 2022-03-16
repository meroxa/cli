package turbinego

import (
	"context"
	"os"
	"os/exec"
	"path"

	"github.com/meroxa/cli/log"
)

// Run will build a go binary and will run it.
func Run(ctx context.Context, appPath string, l log.Logger) error {
	// grab current location to use it as project name
	appName := path.Base(appPath)

	// building is a requirement prior to running for go apps
	err := buildBinary(ctx, l, appPath, appName, false)
	if err != nil {
		return err
	}

	err = os.Chdir(appPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("./" + appName) //nolint:gosec

	out, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, err.Error())
		return err
	}

	l.Info(ctx, string(out))
	return nil
}
