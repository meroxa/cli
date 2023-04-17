package flink

import (
	"context"
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *flinkCLI) Init(_ context.Context, name string) error {
	panic("not implemented")
}

func (t *flinkCLI) GitInit(ctx context.Context, name string) error {
	return utils.GitInit(ctx, name)
}
