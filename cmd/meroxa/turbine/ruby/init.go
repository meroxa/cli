package turbinerb

import (
	"context"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func (t *turbineRbCLI) Init(_ context.Context, appName string) error {
	a := &app.AppInit{
		AppName:  appName,
		Language: ir.Ruby,
		Path:     t.appPath,
	}
	return a.Init()
}
