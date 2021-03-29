package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestDescribeConnectorCmd(t *testing.T) {
	tests := []struct {
		expected string
		args     []string
	}{
		{
			"Error: requires connector name",
			[]string{"describe", "connector"},
		},
		// TODO: Add a test mocking the call when specifying connector name as argument
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
