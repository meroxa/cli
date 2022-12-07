package turbinerb

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/turbine-core/pkg/app"
)

func (t *turbineRbCLI) Init(ctx context.Context, appName string) error {
	return app.NewAppInit(appName, utils.Ruby, t.appPath).Init()
}

func (t *turbineRbCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
