package turbinepy

import (
	"context"
	"os/exec"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbinePyCLI) Run(ctx context.Context) error {
	cmd := exec.Command("turbine-py", "run", t.appPath)
	stdOut, err := utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	t.logger.Info(ctx, stdOut)
	return err
}
