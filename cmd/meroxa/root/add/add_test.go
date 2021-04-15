package add_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/root"
)

func TestAddCmd(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{"Add a resource to your Meroxa resource catalog"},
		{"Usage:\n" +
			"  meroxa add [command]"},
		{"Available Commands:"},
		{"resource    Add a resource to your Meroxa resource catalog"},
		{"Flags:\n" +
			"  -h, --help   help for add"},
	}

	rootCmd := root.Cmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"add"})
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
