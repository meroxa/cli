package turbinepy

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func (t *turbinePyCLI) Init(_ context.Context, appName string) error {
	return app.NewAppInit(appName, ir.Python, t.appPath).Init()
}

func (t *turbinePyCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
