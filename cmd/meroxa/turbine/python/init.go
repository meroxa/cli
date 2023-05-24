package turbinepy

import (
	"context"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func (t *turbinePyCLI) Init(_ context.Context, appName string) error {
	return app.NewAppInit(appName, ir.Python, t.appPath).Init()
}
