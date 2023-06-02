package turbinejava

import (
	"context"
	"fmt"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"os"
	"os/exec"
	"strings"
)

func (t *turbineJavaCLI) Run(ctx context.Context) error {
	gracefulStop, grpcListenAddress := t.Core.Run(ctx)
	defer gracefulStop()

	// TODO: here it goes the command specific in java to run the gRPC server
	//buildCMD := exec.Command("go", []string{
	//	"run", "./...", "build",
	//	"-gitsha", "devel",
	//	"-turbine-core-server", grpcListenAddress,
	//	"-app-path", t.appPath,
	//}...)
	if err := t.build(ctx); err != nil {
		return fmt.Errorf("couldn't build Java app: %w", err)
	}

	return turbine.RunCMD(ctx, t.logger, t.runCMD(ctx, grpcListenAddress))
}

func (t *turbineJavaCLI) build(ctx context.Context) error {
	cmd := exec.Command(
		"mvn",
		"-U",
		"clean",
		"package",
		"-Dquarkus.package.type=uber-jar",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.appPath
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		t.logger.Errorf(ctx, err.Error())
		return err
	}

	return nil
}

func (t *turbineJavaCLI) runCMD(ctx context.Context, grpcAddress string) *exec.Cmd {
	split := strings.Split(grpcAddress, ":")
	cmd := exec.Command(
		"java",
		"-jar",
		"-Dquarkus.grpc.server.port="+split[1],
		"target/turbine-java-example-1.0-SNAPSHOT-runner.jar",
	)

	t.logger.Infof(ctx, "java run cmd: %v", cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.appPath
	cmd.Env = os.Environ()

	return cmd
}
