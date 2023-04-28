package turbinego

import (
	"context"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineGoCLI) Run(ctx context.Context) error {
	go t.runner.Run(ctx)
	defer t.runner.GracefulStop()

	cmd := exec.Command("go", []string{
		"run", "./...", "build",
		"-gitsha", "devel",
		"-turbine-core-server", t.grpcListenAddress,
		"-app-path", t.appPath,
	}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.appPath
	cmd.Env = os.Environ()

	return turbine.RunCMD(ctx, t.logger, cmd)
}
