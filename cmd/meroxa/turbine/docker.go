package turbine

import (
	"context"
	"os"
	"path/filepath"

	"github.com/meroxa/cli/log"
)

type Docker struct{}

func (d *Docker) CleanupDockerfile(logger log.Logger, path string) {
	if err := os.Remove(filepath.Join(path, "Dockerfile")); err != nil {
		logger.Error(context.Background(), err.Error())
	}
}
