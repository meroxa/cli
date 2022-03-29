package turbinepy

import (
	"context"
	"os/exec"

	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func Init(ctx context.Context, l log.Logger, name, path string) error {
	// @TODO sort out binary names so there's no clashing
	cmd := exec.Command("$HOME/.local/bin/turbine", "--generate", name, path)
	_, err := turbineCLI.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}
