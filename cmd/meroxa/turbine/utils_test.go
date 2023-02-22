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
		if tt.path != "" {
			err := os.Mkdir(tt.path, os.ModePerm)
			require.NoError(t, err)
		}

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

func TestCheckGitVersion(t *testing.T) {
	val, err := checkGitVersion(context.Background())
	require.NoError(t, err)
	assert.False(t, val)
}

//nolint:funlen,gocyclo
func TestGitChecks(t *testing.T) {
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

				err = GitInit(ctx, appPath)
				if err != nil {
					return "", err
				}

				// create file
				dockerfilePath := filepath.Join(appPath, "Dockerfile")
				err = os.WriteFile(dockerfilePath, []byte(""), 0644)
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

				err = GitInit(ctx, appPath)
				if err != nil {
					return "", err
				}

				// create file
				dockerfilePath := filepath.Join(appPath, "Dockerfile")
				err = os.WriteFile(dockerfilePath, []byte(""), 0644)
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
				err = os.WriteFile(makefilePath, []byte(""), 0644)
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

				err = GitInit(ctx, appPath)
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

			_, err = GetGitSha(path)
			if err != nil {
				if tc.shaErr == nil {
					t.Fatalf("unepxected error: %v", err)
				}
				assert.Equal(t, tc.shaErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			err = GitChecks(ctx, logger, path)
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

func TestReadAndWriteConfigFile(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()
	appPath, err := makeTmpDir()
	require.NoError(t, err)

	os.Setenv("UNIT_TEST", "true")
	defer func() {
		os.Setenv("UNIT_TEST", "")
	}()

	tests := []struct {
		name  string
		path  string
		input AppConfig
		err   error
	}{
		{
			name: "Successfully read and write AppConfig",
			path: appPath,
			input: AppConfig{
				Name:       "my-name",
				Language:   JavaScript,
				Vendor:     "false",
				ModuleInit: "true",
			},
			err: nil,
		},
		{
			name: "Fail to read and write AppConfig",
			path: "#nope$",
			input: AppConfig{
				Name:       "my-name2",
				Language:   Python3,
				Vendor:     "false",
				ModuleInit: "true",
			},
			err: fmt.Errorf(`open #nope$/app.json: no such file or directory
unable to update app.json file on path "#nope$". Maybe try using a different value for ` + "`" + "--path" + "`"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := WriteConfigFile(tc.path, tc.input)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("unepxected error: %v", err)
				}
				assert.Equal(t, tc.err, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, tc.err)

			lang, err := GetLangFromAppJSON(ctx, logger, tc.path)
			require.NoError(t, err)
			require.Equal(t, tc.input.Language, lang)

			name, err := GetAppNameFromAppJSON(ctx, logger, tc.path)
			require.NoError(t, err)
			require.Equal(t, tc.input.Name, name)

			read, err := ReadConfigFile(tc.path)
			require.NoError(t, err)
			require.Equal(t, tc.input.Vendor, read.Vendor)
			require.Equal(t, tc.input.ModuleInit, read.ModuleInit)

			err = SetVendorInAppJSON(tc.path, true)
			require.NoError(t, err)
			err = SetModuleInitInAppJSON(tc.path, true)
			require.NoError(t, err)

			read, err = ReadConfigFile(tc.path)
			require.NoError(t, err)
			require.Equal(t, "true", read.Vendor)
			require.Equal(t, "false", read.ModuleInit)
		})
	}
}

func TestGetResourceNamesFromString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output []ApplicationResource
	}{
		{
			name:   "Successfully parse names",
			input:  "[one two]",
			output: []ApplicationResource{{Name: "one"}, {Name: "two"}},
		},
		{
			name:   "Successfully parse empty set",
			input:  "[]",
			output: []ApplicationResource{},
		},
		{
			name:   "Successfully parse nonsense",
			input:  "ABC",
			output: []ApplicationResource{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val := GetResourceNamesFromString(tc.input)
			assert.Equal(t, tc.output, val)
		})
	}
}

func TestGetTurbineResponseFromOutput(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
		err    error
	}{
		{
			name:   "Successfully parse output",
			input:  "hey\nturbine-response: very important message\nnot important",
			output: "very important message",
		},
		{
			name:   "Successfully parse empty string",
			input:  "hey\nturbine-response: \nnot important",
			output: "",
		},
		{
			name:  "Fail to find output",
			input: "ABC",
			err:   fmt.Errorf("output is formatted unexpectedly"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, err := GetTurbineResponseFromOutput(tc.input)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("unexpected err: %v", err)
				} else {
					assert.Equal(t, tc.err, err)
				}
			} else {
				assert.Equal(t, tc.output, val)
			}
		})
	}
}

func TestGetPath(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Successfully get path with no input",
			input:  "",
			output: cwd,
		},
		{
			name:   "Successfully get path with . input",
			input:  ".",
			output: cwd,
		},
		{
			name:   "Handle non-existent path",
			input:  "/does+|/not`/exist&",
			output: "/does+|/not`/exist&",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, err := GetPath(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.output, val)
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

func TestRunCMD(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()

	tests := []struct {
		name  string
		input *exec.Cmd
		err   error
	}{
		{
			name:  "Successfully execute command",
			input: exec.Command("date"),
			err:   nil,
		},
		{
			name:  "Fail to find command",
			input: exec.Command("not-a-thing"),
			err:   fmt.Errorf("exec: \"not-a-thing\": executable file not found in $PATH"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := RunCMD(ctx, logger, tc.input)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("unexpected err: %v", err)
				} else {
					assert.Equal(t, tc.err.Error(), err.Error())
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRunCMDWithErrorDetection(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()

	tests := []struct {
		name   string
		input  *exec.Cmd
		output string
		err    error
	}{
		{
			name:  "Successfully execute command",
			input: exec.Command("date"),
			err:   nil,
		},
		{
			name:  "Fail to find command",
			input: exec.Command("not-a-thing"),
			err:   fmt.Errorf("exec: \"not-a-thing\": executable file not found in $PATH"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, err := RunCmdWithErrorDetection(ctx, tc.input, logger)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("unexpected err: %v", err)
				} else {
					assert.Equal(t, tc.err, err)
					assert.Equal(t, tc.output, val)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_trimNonNpmErrorLines(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Successfully remove errors",
			input:  "hi\nnpm warn secrets\nnpm error no good\nmessage",
			output: "hi\nnpm error no good\nmessage",
		},
		{
			name:   "Nothing to remove",
			input:  "hi\nmessage",
			output: "hi\nmessage",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val := trimNonNpmErrorLines(tc.input)
			assert.Equal(t, tc.output, val)
		})
	}
}

func TestCleanupDockerfile(t *testing.T) {
	logger := log.NewTestLogger()

	appPath, err := makeTmpDir()
	require.NoError(t, err)

	dockerfilePath := filepath.Join(appPath, "Dockerfile")
	err = os.WriteFile(dockerfilePath, []byte(""), 0644)
	require.NoError(t, err)

	_, err = os.Stat(dockerfilePath)
	require.NoError(t, err)

	CleanupDockerfile(logger, appPath)
	_, err = os.Stat(dockerfilePath)
	require.True(t, os.IsNotExist(err))
}

func makeTmpDir() (string, error) {
	basePath := "/tmp"
	dirName := uuid.NewString()
	appPath := filepath.Join(basePath, dirName)
	err := os.MkdirAll(appPath, os.ModePerm)
	return appPath, err
}
