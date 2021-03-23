package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestUpdateResourceCmd(t *testing.T) {
	tests := []struct {
		expected string
		args []string
	}{
		{
			"Error: requires a resource name and either `--metadata`, `--url` or `--credentials` to update the resource",
			[]string{"update", "resource"},
		},
		{
			"Error: requires a resource name and either `--metadata`, `--url` or `--credentials` to update the resource",
			[]string{"update", "resource", "--metadata", "{\"logical_replication\": \"true\"}"},
		},
		{
			"Error: requires a resource name and either `--metadata`, `--url` or `--credentials` to update the resource",
			[]string{"update", "resource", "name"},
		},
		// TODO: Add a test mocking the call when specifying resource name and one of the required flags
	}

	for _, tt := range tests {
		rootCmd := RootCmd()
		b := bytes.NewBufferString("")
		rootCmd.SetOut(b)
		rootCmd.SetErr(b)
		rootCmd.SetArgs(tt.args)
		rootCmd.Execute()
		output, err := ioutil.ReadAll(b)

		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(output), tt.expected) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.expected, string(output))
		}
	}
}
