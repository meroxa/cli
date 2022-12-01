package turbinejs

import (
	"context"
)

// Build has nothing to do for turbine-js.
func (t *turbineJsCLI) Build(ctx context.Context, appName string, platform bool) error {
	return nil
}

func (t *turbineJsCLI) CleanUpTempBuildLocation(ctx context.Context) error {
	return nil
}

func (t *turbineJsCLI) CleanUpBinaries(ctx context.Context, appName string) {
}
