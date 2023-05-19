package turbinerb

import (
	"context"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func (t *turbineRbCLI) Init(_ context.Context, appName string) error {
	return app.NewAppInit(appName, ir.Ruby, t.appPath).Init()
}
