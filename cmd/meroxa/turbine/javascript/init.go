package turbinejs

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

func Init(ctx context.Context, l log.Logger, name, path string) error {
	cmd := RunTurbineJS("generate", name, path)
	_, err := utils.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}
