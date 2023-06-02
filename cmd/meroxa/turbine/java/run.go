package turbinejava

import (
	"context"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineJavaCLI) Run(ctx context.Context) error {
	gracefulStop, grpcListenAddress := t.Core.Run(ctx)
	defer gracefulStop()

	// TODO: here it goes the command specific in java to run the gRPC server
	//cmd := exec.Command("go", []string{
	//	"run", "./...", "build",
	//	"-gitsha", "devel",
	//	"-turbine-core-server", grpcListenAddress,
	//	"-app-path", t.appPath,
	//}...)

	cmd := exec.Command("", grpcListenAddress)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.appPath
	cmd.Env = os.Environ()

	return turbine.RunCMD(ctx, t.logger, cmd)
}
