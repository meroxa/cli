package deprecated_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/root"
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

	rootCmd := root.Cmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"describe"})
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
