package turbinerb

import (
	"context"
	"errors"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineRbCLI) Init(ctx context.Context, name string) error {
	return errors.New("unimplemented")
}

func (t *turbineRbCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
