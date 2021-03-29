package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestUpdatePipelineCmd(t *testing.T) {
	tests := []struct {
		expected string
		args     []string
	}{
		{
			"Error: requires pipeline name",
			[]string{"update", "pipeline"},
		},
		{
			"Error: required flag(s) \"state\" not set",
			[]string{"update", "pipeline", "name"},
		},
		// TODO: Add a test mocking the call when specifying --state
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
