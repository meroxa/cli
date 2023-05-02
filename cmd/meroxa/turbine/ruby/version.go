package turbinerb

import (
	"context"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
)

// GetVersion will return the version of turbine-rb dependency of a given app.
func (t *turbineRbCLI) GetVersion(ctx context.Context) (string, error) {
	ver, err := internal.RunTurbineCmd(
		ctx,
		t.appPath,
		internal.TurbineCommandVersion,
		map[string]string{},
	)
	if err != nil {
		return "", fmt.Errorf("cannot determine turbine-rb gem version in %s", t.appPath)
	}
	return ver, nil
}
