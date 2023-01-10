package turbinego

import (
	"context"
	"os"
	"path/filepath"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

// Run will build a go binary and will run it.
func (t *turbineGoCLI) Run(ctx context.Context) error {
	// building is a requirement prior to running for go apps
	if err := t.Build(ctx, t.appName, false); err != nil {
		return err
	}

	go t.runServer.Run(ctx)
	defer t.runServer.GracefulStop()

	cmd := NewTurbineGoCmd(t.appPath, t.appName, TurbineCommandRun, map[string]string{
		"TURBINE_CORE_SERVER": t.grpcListenAddress,
	})
	return turbine.RunCMD(ctx, t.logger, cmd)
}

// RunCleanup removes any dangling binaries.
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
