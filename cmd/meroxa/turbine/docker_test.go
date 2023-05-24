package turbine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meroxa/cli/log"
	"github.com/stretchr/testify/require"
)

func TestCleanupDockerfile(t *testing.T) {
	var d Docker

	logger := log.NewTestLogger()

	appPath, err := makeTmpDir()
	require.NoError(t, err)

	dockerfilePath := filepath.Join(appPath, "Dockerfile")
	err = os.WriteFile(dockerfilePath, []byte(""), 0o644)
	require.NoError(t, err)

	_, err = os.Stat(dockerfilePath)
	require.NoError(t, err)

	d.CleanupDockerfile(logger, appPath)
	_, err = os.Stat(dockerfilePath)
	require.True(t, os.IsNotExist(err))
}
