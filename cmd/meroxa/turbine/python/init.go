package turbinepy

import (
	"context"
	"os/exec"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbinePyCLI) Init(ctx context.Context, name string) error {
	cmd := exec.Command("turbine-py", "generate", name, t.appPath)
	_, err := utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	return err
}

func (t *turbinePyCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
