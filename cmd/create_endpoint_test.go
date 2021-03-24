package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCreateEndpointCmd(t *testing.T) {
	tests := []struct {
		expected string
		args []string
	}{
		{
			"Error: required flag(s) \"protocol\", \"stream\" not set",
			[]string{"create", "endpoint"},
		},
		{
			"Error: required flag(s) \"stream\" not set",
			[]string{"create", "endpoint", "--protocol", "http"},
		},
		{
			"Error: required flag(s) \"protocol\" not set",
			[]string{"create", "endpoint", "--stream", "my-strea,"},
		},
		// TODO: Add a test with "--protocol" and "--stream", mocking the call
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
