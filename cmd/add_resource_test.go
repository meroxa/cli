package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestAddResourceCmd(t *testing.T) {
	tests := []struct {
		expected string
		args []string
	}{
		{
			"Error: required flag(s) \"type\", \"url\" not set",
			[]string{"add", "resource"},
		},
		{
			"Error: required flag(s) \"type\" not set",
			[]string{"add", "resource", "--url", "myUrl"},
		},
		{
			"Error: requires resource name",
			[]string{"add", "resource", "--url", "myUrl", "--type", "postgres"},
		},
		// TODO: Add a test with resource name as argument and mocking the call
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
