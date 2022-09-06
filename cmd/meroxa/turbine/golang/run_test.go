package turbinego

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/meroxa/cli/log"
)

func TestRunCleanup(t *testing.T) {
	logger := log.NewTestLogger()
	ctx := context.Background()

	tests := []struct {
		description string
		appName     string
		binaries    []string
		output      bool
	}{
		{
			description: "Successfully delete binaries",
			appName:     "verycool",
			binaries:    []string{"verycool", "verycool.cross"},
			output:      false,
		},
		{
			description: "Successfully delete binaries",
			appName:     "changed",
			binaries:    []string{"verycool", "verycool.cross"},
			output:      true,
		},
	}

	cwd, _ := filepath.Abs(".")
	defer func() {
		RunCleanup(ctx, logger, cwd, "verycool")
	}()

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			for _, name := range tc.binaries {
				err := os.WriteFile(name, []byte{}, 0644)
				require.NoError(t, err)
			}

			RunCleanup(ctx, logger, cwd, tc.appName)
			output := logger.LeveledOutput()
			if tc.output {
				require.Greater(t, len(output), 0)
			} else {
				require.Equal(t, len(output), 0)
			}
		})
	}
}
