package old_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/root"
)

func TestUpdateCmd(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{"Update a component of the Meroxa platform, including connectors"},
		{"Usage:\n" +
			"  meroxa update [command]"},
		{"Available Commands:"},
		{"connector   Update connector state"},
		{"pipeline    Update pipeline state"},
		{"resource    Update a resource"},
		{"Flags:\n" +
			"  -h, --help   help for update"},
	}

	rootCmd := root.Cmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"update"})
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
