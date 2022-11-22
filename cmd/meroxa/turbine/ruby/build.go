package turbinerb

import (
	"context"
	"fmt"
)

func (t *turbineRbCLI) Build(ctx context.Context, appName string, platform bool) (string, error) {
	return "", fmt.Errorf("build command not implemented")
}

func (t *turbineRbCLI) CleanUpTempBuildLocation(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (t *turbineRbCLI) CleanUpBinaries(ctx context.Context, appName string) {
}
