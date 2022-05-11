package turbinejs

import (
	"context"

	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

func Init(ctx context.Context, l log.Logger, name, path string) error {
	cmd := turbinecli.RunTurbineJS("generate", name, path)
	_, err := turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}
