package turbinepy

import (
	"context"
	"os/exec"
	"path/filepath"

	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func Run(ctx context.Context, l log.Logger, path string) error {
	f, err := filepath.Abs(filepath.Join(path, "main.py"))
	if err != nil {
		return err
	}
	cmd := exec.Command("python3", f)
	_, err = turbineCLI.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}
