package deprecated_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/root"
)

func TestUpdateConnectorCmd(t *testing.T) {
	tests := []struct {
		expected string
		args     []string
	}{
		{
			"Error: requires connector name",
			[]string{"update", "connector"},
		},
		{
			"Error: required flag(s) \"state\" not set",
			[]string{"update", "connector", "name"},
		},
		// TODO: Add a test mocking the call when specifying --state
	}

	for _, tt := range tests {
		rootCmd := root.Cmd()
		b := bytes.NewBufferString("")
		rootCmd.SetOut(b)
		rootCmd.SetErr(b)
		rootCmd.SetArgs(tt.args)
		_ = rootCmd.Execute()
		output, err := ioutil.ReadAll(b)

		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(output), tt.expected) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.expected, string(output))
		}
	}
}
