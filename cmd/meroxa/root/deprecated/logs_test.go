package deprecated_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/root"
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

	rootCmd := root.Cmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"logs"})
	_ = rootCmd.Execute()

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
