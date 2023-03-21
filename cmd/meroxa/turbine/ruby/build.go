package turbinerb

import (
	"context"
)

func (t *turbineRbCLI) Build(_ context.Context, _ string, _ bool) error {
	return nil
}

func (t *turbineRbCLI) CleanUpBinaries(_ context.Context, _ string) {
}
