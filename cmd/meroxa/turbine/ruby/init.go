package turbinerb

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineRbCLI) Init(ctx context.Context, name string) error {
	return nil
}

func (t *turbineRbCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
