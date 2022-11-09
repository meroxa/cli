package turbinerb

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/turbine/core"
)

const GrpcServerPort = 50500

func (t *turbineRbCLI) Run(ctx context.Context) (err error) {
	// Run turbine-core gRPC server
	cs := core.NewTurbineCoreServer()
	go cs.Run(GrpcServerPort, "local")

	// Execute Turbine.run on app:
	// turbine_client process
	pipelineCmd := exec.Command("turbine_client", "process")
	pipelineCmd.Env = append(os.Environ(),
		fmt.Sprintf("TURBINE_CORE_SERVER=localhost:%d", GrpcServerPort))
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
	err = pipelineCmd.Wait()
	if err != nil {
		t.logger.Errorf(ctx, err.Error())
		return err
	}

	// teardown turbine core server
	cs.Stop()

	return err
}
