package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCreateConnectorCmd(t *testing.T) {
	tests := []struct {
		expected string
		args     []string
	}{
		{
			"Error: requires either a source (--from) or a destination (--to)",
			[]string{"create", "connector"},
		},
		{
			"Error: required flag(s) \"input\" not set",
			[]string{"create", "connector", "--to", "pg2redshift"},
		},
		{
			"Error: required flag(s) \"input\" not set",
			[]string{"create", "connector", "--from", "pg2kafka"},
		},
		// TODO: Add a test with "--input" and mocking the call
		// TODO: Add a test with connector name as argument and mocking the call
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
