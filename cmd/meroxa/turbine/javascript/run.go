package turbinejs

import "context"

func (j *turbineJsCLI) Run(ctx context.Context) error {
	return Build(ctx, j.logger, j.appPath)
}
