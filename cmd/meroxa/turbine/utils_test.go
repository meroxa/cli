package turbine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadAndWriteConfigFile(t *testing.T) {
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
				Language:   ir.JavaScript,
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
				Language:   ir.Python,
				Vendor:     "false",
				ModuleInit: "true",
			},
			err: fmt.Errorf(`open #nope$/app.json: no such file or directory
unable to update app.json file on path "#nope$". Maybe try using a different value for ` + "`" + "--path" + "`"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := WriteConfigFile(tc.path, &tc.input)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				assert.Equal(t, tc.err, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, tc.err)

			lang, err := GetLangFromAppJSON(logger, tc.path)
			require.NoError(t, err)
			require.Equal(t, tc.input.Language, lang)

			name, err := GetAppNameFromAppJSON(logger, tc.path)
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
			name:   "Successfully parse with carriage returns",
			input:  "hey\nturbine-response: 1.7.0\r\nthis is from windows!",
			output: "1.7.0",
		},
		{
			name:   "Successfully parse with carriage returns plus some spaces",
			input:  "hey\nturbine-response: 1.7.0   \r\nthis is from windows!",
			output: "1.7.0",
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

func makeTmpDir() (string, error) {
	basePath := "/tmp"
	dirName := uuid.NewString()
	appPath := filepath.Join(basePath, dirName)
	err := os.MkdirAll(appPath, os.ModePerm)
	return appPath, err
}
