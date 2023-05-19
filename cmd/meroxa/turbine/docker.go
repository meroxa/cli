package turbine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/meroxa/cli/log"
)

type Docker struct{}

func (d *Docker) CleanupDockerfile(logger log.Logger, path string) {
	logger.StartSpinner("\t", fmt.Sprintf("Removing Dockerfile created for your application in %s...", path))
	// We clean up Dockerfile as last step

	if err := os.Remove(filepath.Join(path, "Dockerfile")); err != nil {
		logger.StopSpinnerWithStatus(fmt.Sprintf("Unable to remove Dockerfile: %v", err), log.Failed)
	} else {
		logger.StopSpinnerWithStatus("Dockerfile removed", log.Successful)
	}
}
