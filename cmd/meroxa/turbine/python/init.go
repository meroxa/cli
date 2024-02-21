package turbinepy

import (
	"context"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func (t *turbinePyCLI) Init(_ context.Context, appName string) error {
	a := &app.AppInit{
		AppName:  appName,
		Language: ir.Python,
		Path:     t.appPath,
	}
	return a.Init()
}
