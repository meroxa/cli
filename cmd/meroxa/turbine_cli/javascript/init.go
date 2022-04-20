package turbinejs

import (
	"context"
	"os/exec"

	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func Init(ctx context.Context, l log.Logger, name, path string) error {
	cmd := exec.Command("npx", "--yes", "@meroxa/turbine-js@0.1.4", "generate", name, path)
	_, err := turbineCLI.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}
