package turbinerb

import (
	"context"
)

// Build has nothing to do for turbine-rb.
func (t *turbineRbCLI) Build(ctx context.Context, appName string, platform bool) (string, error) {
	return "", nil
}

func (t *turbineRbCLI) CleanUpTempBuildLocation(ctx context.Context) error {
	return nil
}

func (t *turbineRbCLI) CleanUpBinaries(ctx context.Context, appName string) {
}
