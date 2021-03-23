package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestListCmd(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{"List the components of the Meroxa platform, including pipelines,\n" +
			" resources, connectors, etc... You may also filter by type."},
		{"Usage:\n  meroxa list [command]\n\n"},
		{"Available Commands:"},
		{"connectors     List connectors"},
		{"endpoint       List endpoints"},
		{"pipelines      List pipelines"},
		{"resource-types List resources-types"},
		{"resources      List resources"},
		{"transforms     List transforms"},
		{"Flags:\n  -h, --help   help for list\n"},
	}

	rootCmd := RootCmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"list"})
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
