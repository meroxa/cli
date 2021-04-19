package old_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/root"
)

func TestCreateCmd(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{"Use the create command to create various Meroxa pipeline components\n" +
			"including connectors."},
		{"Usage:\n  meroxa create [command]"},
		{"Available Commands:"},
		{"connector   Create a connector"},
		{"endpoint    Create an endpoint"},
		{"pipeline    Create a pipeline"},
		{"Flags:\n  -h, --help   help for create\n"},
	}

	rootCmd := root.Cmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"create"})
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
