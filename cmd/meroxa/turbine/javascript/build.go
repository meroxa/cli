package turbinejs

import (
	"context"
)

// Build has nothing to do for turbine-js.
func (t *turbineJsCLI) Build(_ context.Context, _ string, _ bool) error {
	return nil
}

func (t *turbineJsCLI) CleanUpTempBuildLocation(_ context.Context) error {
	return nil
}

func (t *turbineJsCLI) CleanUpBinaries(_ context.Context, _ string) {
}
