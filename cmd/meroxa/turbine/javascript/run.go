package turbinejs

import (
	"context"
)

// Run calls turbine-js to exercise the app.
func (t *turbineJsCLI) Run(ctx context.Context) (err error) {
	stdOut, err := RunTurbineJS(ctx, t.logger, "test", t.appPath)
	t.logger.Info(ctx, stdOut)
	return err
}
