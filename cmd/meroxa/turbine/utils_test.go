package turbine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGitInit(t *testing.T) {
	testDir := os.TempDir() + "/tests" + uuid.New().String()

	tests := []struct {
		path string
		err  error
	}{
		{path: "", err: errors.New("path is required")},
		{path: testDir, err: nil},
	}

	for _, tt := range tests {
		err := GitInit(context.Background(), tt.path)
		if err != nil {
			if tt.err == nil {
				t.Fatalf("unexpected error \"%s\"", err)
			} else if tt.err.Error() != err.Error() {
				t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
			}
		}

		if tt.err == nil {
			if _, err := os.Stat(testDir + "/.git"); os.IsNotExist(err) {
				t.Fatalf("expected directory \"%s\" to be created", testDir)
			}
		}
	}

	os.RemoveAll(testDir)
}

func TestGetTurbineJSBinary(t *testing.T) {
	testCases := []struct {
		name    string
		envVar  string
		wantCmd string
	}{
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is unset",
			envVar:  "",
			wantCmd: fmt.Sprintf("@meroxa/turbine-js-cli@%s", turbineJSVersion),
		},
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is true",
			envVar:  "true",
			wantCmd: "turbine-js-cli",
		},
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is false",
			envVar:  "false",
			wantCmd: fmt.Sprintf("@meroxa/turbine-js-cli@%s", turbineJSVersion),
		},
		{
			name:    "MEROXA_USE_LOCAL_TURBINE_JS is set to a value that is neither true nor false",
			envVar:  "jam",
			wantCmd: fmt.Sprintf("@meroxa/turbine-js-cli@%s", turbineJSVersion),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("MEROXA_USE_LOCAL_TURBINE_JS", tc.envVar)

			params := []string{"foo", "bar"}
			result := getTurbineJSBinary(params)
			assert.Equal(t, []string{"npx", "--yes", tc.wantCmd, "foo", "bar"}, result)
		})
	}
}
