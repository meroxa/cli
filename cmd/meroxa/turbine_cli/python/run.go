package turbinepy

import (
	"context"
	"os/exec"

	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func Run(ctx context.Context, l log.Logger, path string) error {
	cmd := exec.Command("turbine-py", "run", path)
	stdOut, err := turbineCLI.RunCmdWithErrorDetection(ctx, cmd, l)
	l.Info(ctx, stdOut)
	return err
}
