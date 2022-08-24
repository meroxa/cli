package turbinejs

import "context"

func (t *turbineJsCLI) Run(ctx context.Context) error {
	return Build(ctx, t.logger, t.appPath)
}
