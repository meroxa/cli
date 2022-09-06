package turbinejs

import "context"

func (t *turbineJsCLI) Run(ctx context.Context) (err error) {
	_, err = t.Build(ctx, "", false)
	return err
}
