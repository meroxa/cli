package turbinerb

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/turbine/core"
)

const GRPC_SERVER_PORT = 50500

func (t *turbineRbCLI) Run(ctx context.Context) error {
	// Run turbine-core gRPC server
	cs := core.NewCoreServer()
	go cs.Run(GRPC_SERVER_PORT, "local")

	// Execute Turbine.run on app:
	// bundle exec ruby -e "require_relative('./app'); Turbine.run;"
	pipelineCmd := exec.Command("bundle", "exec", "ruby", "-r", "./app", "-e", "Turbine.run")
	pipelineCmd.Env = append(os.Environ(),
		fmt.Sprintf("TURBINE_CORE_SERVER=localhost:%d", GRPC_SERVER_PORT))
	pipelineCmd.Stdout = os.Stdout
	pipelineCmd.Stderr = os.Stderr

	// set working directory to where the CLI was run
	folderPath, err := os.Getwd()
	if err != nil {
		t.logger.Error(ctx, err.Error())
		return err
	}
	pipelineCmd.Dir = folderPath

	// Run command
	err = pipelineCmd.Start()
	if err != nil {
		t.logger.Errorf(ctx, err.Error())
		return err
	}

	// wait for pipelineCmd to exit
	pipelineCmd.Wait()

	// teardown turbine core server
	cs.Stop()

	return err
}
