package flink

import (
	"context"
)

func (t *flinkCLI) Build(_ context.Context, _ string, _ bool) error {
	return nil
}

func (t *flinkCLI) CleanUpBinaries(_ context.Context, _ string) {
}
