package turbine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/meroxa/cli/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitInit(t *testing.T) {
	var g Git
	testDir := os.TempDir() + "/tests" + uuid.New().String()

	tests := []struct {
		path string
		err  error
	}{
		{path: "", err: errors.New("path is required")},
		{path: testDir, err: nil},
	}

	for _, tt := range tests {
		if tt.path != "" {
			err := os.Mkdir(tt.path, os.ModePerm)
			require.NoError(t, err)
		}

		err := g.GitInit(context.Background(), tt.path)
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

func TestCheckGitVersion(t *testing.T) {
	val, err := checkGitVersion(context.Background())
	require.NoError(t, err)
	assert.False(t, val)
}

//nolint:funlen,gocyclo
func TestCheckUncommittedChanges(t *testing.T) {
	var g Git

	if gh := os.Getenv("GITHUB_WORKSPACE"); gh != "" {
		t.Skipf("skipping git test in github action")
	}

	ctx := context.Background()
	logger := log.NewTestLogger()

	tests := []struct {
		desc      string
		setup     func() (string, error)
		branch    string
		shaErr    error
		checksErr error
	}{
		{
			desc: "Successfully pass all git commands",
			setup: func() (string, error) {
				appPath, err := makeTmpDir()
				if err != nil {
					return "", err
				}

				err = g.GitInit(ctx, appPath)
				if err != nil {
					return "", err
				}

				// create file
				dockerfilePath := filepath.Join(appPath, "Dockerfile")
				err = os.WriteFile(dockerfilePath, []byte(""), 0o644)
				if err != nil {
					return "", err
				}

				// add file
				cmd := exec.Command("git", "add", "Dockerfile")
				cmd.Dir = appPath
				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("%s", string(output))
				}

				// commit file
				cmd = exec.Command("git", "commit", "-m", "first")
				cmd.Dir = appPath
				output, err = cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("%s", string(output))
				}

				return appPath, nil
			},
			branch: "main",
		},
		{
			desc: "Fail uncommitted changes",
			setup: func() (string, error) {
				appPath, err := makeTmpDir()
				if err != nil {
					return "", err
				}

				err = g.GitInit(ctx, appPath)
				if err != nil {
					return "", err
				}

				// create file
				dockerfilePath := filepath.Join(appPath, "Dockerfile")
				err = os.WriteFile(dockerfilePath, []byte(""), 0o644)
				if err != nil {
					return "", err
				}

				// feature branch
				cmd := exec.Command("git", "checkout", "-b", "unit-test")
				cmd.Dir = appPath
				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("%s", string(output))
				}

				// add file
				cmd = exec.Command("git", "add", "Dockerfile")
				cmd.Dir = appPath
				output, err = cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("%s", string(output))
				}

				// commit file
				cmd = exec.Command("git", "commit", "-m", "first")
				cmd.Dir = appPath
				output, err = cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("%s", string(output))
				}

				// create second file
				makefilePath := filepath.Join(appPath, "Makefile")
				err = os.WriteFile(makefilePath, []byte(""), 0o644)
				if err != nil {
					return "", err
				}

				// add second file
				cmd = exec.Command("git", "add", "Makefile")
				cmd.Dir = appPath
				output, err = cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("%s", string(output))
				}

				return appPath, nil
			},
			branch:    "unit-test",
			checksErr: fmt.Errorf("unable to proceed with deployment because of uncommitted changes"),
		},
		{
			desc: "Fail git SHA check",
			setup: func() (string, error) {
				appPath, err := makeTmpDir()
				if err != nil {
					return "", err
				}

				err = g.GitInit(ctx, appPath)
				if err != nil {
					return "", err
				}

				return appPath, nil
			},
			branch: "main",
			//nolint:revive
			shaErr: fmt.Errorf(
				`/usr/bin/git rev-parse HEAD: fatal: ambiguous argument 'HEAD': unknown revision or path not in the working tree.
Use '--' to separate paths from revisions, like this:
'git <command> [<revision>...] -- [<file>...]'
HEAD
`),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			path, err := tc.setup()
			require.NoError(t, err)

			branch, err := GetGitBranch(path)
			require.NoError(t, err)
			assert.Equal(t, tc.branch, branch)

			_, err = g.GetGitSha(ctx, path)
			if err != nil {
				if tc.shaErr == nil {
					t.Fatalf("unepxected error: %v", err)
				}
				assert.Equal(t, tc.shaErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			err = g.CheckUncommittedChanges(ctx, logger, path)
			if err != nil {
				if tc.checksErr == nil {
					t.Fatalf("unepxected error: %v", err)
				}
				assert.Equal(t, tc.checksErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGitMainBranch(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output bool
	}{
		{
			name:   "Process main branch",
			input:  "main",
			output: true,
		},
		{
			name:   "Process master branch",
			input:  "master",
			output: true,
		},
		{
			name:   "Process other branch",
			input:  "anything-else",
			output: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val := GitMainBranch(tc.input)
			assert.Equal(t, tc.output, val)
		})
	}
}
