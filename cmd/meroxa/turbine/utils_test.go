package turbine

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meroxa/cli/log"
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

func TestUploadFile(t *testing.T) {
	ctx := context.Background()
	retries := 0
	testCases := []struct {
		name    string
		server  func(int) *httptest.Server
		status  int
		retries int
		output  string
		err     error
	}{
		{
			name: "Successfully upload file",
			server: func(status int) *httptest.Server {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					retries++
					w.WriteHeader(status)
				}))
				return server
			},
			status:  http.StatusOK,
			retries: 1,
			output:  "Source uploaded",
		},
		{
			name: "Fail to upload file",
			server: func(status int) *httptest.Server {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					retries++
					w.WriteHeader(status)
				}))
				return server
			},
			status:  http.StatusInternalServerError,
			retries: 3,
			output:  "Failed to upload build source file",
			err:     fmt.Errorf("upload failed: 500 Internal Server Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			retries = 0
			logger := log.NewTestLogger()
			logger.StartSpinner("", "")
			server := tc.server(tc.status)
			err := uploadFile(ctx, logger, "utils.go", server.URL)
			if tc.err != nil {
				assert.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.retries, retries)
			output := logger.SpinnerOutput()
			assert.True(t, strings.Contains(output, tc.output))
			server.Close()
		})
	}
}
