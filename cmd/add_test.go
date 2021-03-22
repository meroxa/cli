package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
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

	rootCmd := RootCmd()
	listCmd := ListCmd()
	rootCmd.AddCommand(listCmd)

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
