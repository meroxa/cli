package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestDescribeCmd(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{"Describe a component of the Meroxa data platform, including resources and connectors"},
		{"Usage:\n" +
			"  meroxa describe [command]"},
		{"Available Commands:"},
		{"connector   Describe connector"},
		{"endpoint    Describe endpoint"},
		{"resource    Describe resource"},
		{"Flags:\n" +
			"  -h, --help   help for describe"},
	}

	rootCmd := RootCmd()
	listCmd := ListCmd()
	rootCmd.AddCommand(listCmd)

	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"describe"})
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
