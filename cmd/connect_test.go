package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestConnectCmd(t *testing.T) {
	tests := []struct {
		expected string
		args []string
	}{
		{
			"Error: required flag(s) \"from\", \"to\" not set",
			[]string{"connect"},
		},
		{
			"Error: required flag(s) \"to\" not set",
			[]string{"connect", "--from", "resource-name"},
		},
		{
			"Error: required flag(s) \"from\" not set",
			[]string{"connect", "--to", "resource-name"},
		},
		// TODO: Add a test with connect --to and --from mocking the call
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
