package turbinejs

import (
	"context"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/turbine/javascript/internal"
)

func (t *turbineJsCLI) GetVersion(ctx context.Context) (string, error) {
	ver, err := internal.RunTurbineCmd(
		ctx,
		t.appPath,
		internal.TurbineCommandVersion,
		map[string]string{},
	)
	if err != nil {
		return "", fmt.Errorf("cannot determine turbine-js module version in %s", t.appPath)
	}
	return ver, nil
}
