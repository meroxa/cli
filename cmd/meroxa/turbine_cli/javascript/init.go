package turbinejs

import (
	"context"
	"os/exec"

	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func Init(ctx context.Context, l log.Logger, name, path string) error {
	cmd := exec.Command("npx", "turbine", "generate", name, path)
	_, err := turbineCLI.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}
