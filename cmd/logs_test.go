package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestLogsCmd(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{"Print logs for a component"},
		{"Usage:\n  meroxa logs [command]\n\n"},
		{"Available Commands:"},
		{"connector   Print logs for a connector"},
		{"Flags:\n  -h, --help   help for logs\n"},
	}

	rootCmd := RootCmd()
	listCmd := ListCmd()
	rootCmd.AddCommand(listCmd)

	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"logs"})
	rootCmd.Execute()

	out, err := ioutil.ReadAll(b)

	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		if !strings.Contains(string(out), tt.expected) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.expected, string(out))
		}
	}
}
