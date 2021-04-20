package resource_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/root"
)

func TestResourceCmd(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{"Manage resources on Meroxa"},
		{"Usage:\n" +
			"  meroxa resources [command]"},
		{"Available Commands:"},
		{"create      Create a resource into your Meroxa resource catalog"},
		{"Flags:\n" +
			"  -h, --help   help for resources"},
	}

	rootCmd := root.Cmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"resources"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		if !strings.Contains(string(out), tt.expected) {
			t.Fatalf("expected %q got %q", tt.expected, string(out))
		}
	}
}
