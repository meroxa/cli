package turbinepy

import (
	"context"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/turbine/python/internal"
)

// GetVersion calls turbine-py to get its version used by the given app.
func (t *turbinePyCLI) GetVersion(_ context.Context) (string, error) {
	ver, err := internal.RunTurbineCmd(
		t.appPath,
		internal.TurbineCommandVersion,
		map[string]string{},
		t.appPath,
	)
	if err != nil {
		return "", fmt.Errorf("cannot determine turbine-py version in %s", t.appPath)
	}
	return ver, nil
}
