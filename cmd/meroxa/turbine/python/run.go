package turbinepy

import (
	"context"
	"os/exec"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (p *turbinePyCLI) Run(ctx context.Context) error {
	cmd := exec.Command("turbine-py", "run", p.appPath)
	stdOut, err := utils.RunCmdWithErrorDetection(ctx, cmd, p.logger)
	p.logger.Info(ctx, stdOut)
	return err
}
