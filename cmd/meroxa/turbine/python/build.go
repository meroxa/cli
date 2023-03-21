package turbinepy

import (
	"context"
)

// Build created the needed structure for a python app.
func (t *turbinePyCLI) Build(_ context.Context, _ string, _ bool) error {
	return nil
}

func (t *turbinePyCLI) CleanUpBinaries(_ context.Context, _ string) {
}
