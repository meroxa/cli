package turbinerb

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/turbine/core"
)

func (t *turbineRbCLI) Run(ctx context.Context) error {
	//todo:
	//- launch turbine-core grpc server
	//- execute pipeline_runner.sh
	//- kill turbine-core grpc server
	t.logger.Info(ctx, "running Ruby app")

	// Run turbine-core gRPC server
	cs := core.NewCoreServer()
	go cs.Run(50500, "local")

	// Execute Turbine.run on app:
	// bundle exec ruby -e "require_relative('./app'); Turbine.run;"
	pipelineCmd := exec.Command("bundle", "exec", "ruby", "-e", "require_relative('./app');", "Turbine.run;")
	t.logger.Infof(ctx, "dir: %s", pipelineCmd.Dir)

	// set working directory to where the CLI was run
	folderPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		t.logger.Error(ctx, err.Error())
	}
	pipelineCmd.Dir = folderPath
	t.logger.Infof(ctx, "dir: %s", pipelineCmd.Dir)
	//pipelineCmd := exec.Command("bundle", "exec", "ruby", "-e", "require_relative('./app')")
	//pipelineCmd := exec.Command("sleep", "5")
	stdout, stderr := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	pipelineCmd.Stdout = stdout
	pipelineCmd.Stderr = stderr
	t.logger.Info(ctx, "starting pipeline cmd")

	err = pipelineCmd.Start()
	if err != nil {
		t.logger.Errorf(ctx, err.Error())
		return err
	}
	stdOutMsg := stdout.String()
	stdErrMsg := stderr.String()
	t.logger.Info(ctx, stdOutMsg)
	if stdErrMsg != "" {
		t.logger.Error(ctx, stdErrMsg)
	}

	// wait for pipelineCmd to exit
	t.logger.Info(ctx, "waiting on pipelineCmd")
	pipelineCmd.Wait()
	t.logger.Info(ctx, "done waiting")

	// teardown turbine core server
	t.logger.Info(ctx, "killing turbine-core")
	cs.Stop()
	t.logger.Info(ctx, "killed turbine-core")

	return err
}
