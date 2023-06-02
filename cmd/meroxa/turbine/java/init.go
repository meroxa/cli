package turbinejava

import (
	"context"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func (t *turbineJavaCLI) Init(_ context.Context, appName string) error {
	return app.NewAppInit(appName, ir.Java, t.appPath).Init()
}
