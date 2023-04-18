package turbinerb

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func (t *turbineRbCLI) Init(_ context.Context, appName string) error {
	return app.NewAppInit(appName, ir.Ruby, t.appPath).Init()
}

func (t *turbineRbCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
