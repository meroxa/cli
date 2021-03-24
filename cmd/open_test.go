package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestOpenCmd(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{"Open in a web browser"},
		{"Usage:\n" +
			"  meroxa open [command]"},
		{"Available Commands:"},
		{"billing     Open your billing page in a web browser"},
		{"Flags:\n" +
			"  -h, --help   help for open"},
	}

	rootCmd := RootCmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"open"})
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
