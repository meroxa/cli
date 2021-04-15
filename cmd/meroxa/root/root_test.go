package root

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCmd(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{"The Meroxa CLI allows quick and easy access to the Meroxa data platform.\n\n" +
			"Using the CLI you are able to create and manage sophisticated data pipelines\n" +
			"with only a few simple commands. You can get started by listing the supported\n" +
			"resource types:\n\n" +
			"meroxa list resource-types\n\n"},
		{"Usage:\n  meroxa [command]\n\n"},
		{"Available Commands:"},
		{"add         Add a resource to your Meroxa resource catalog"},
		{"api         Invoke Meroxa API"},
		{"billing     Open your billing page in a web browser"},
		{"completion  Generate completion script"},
		{"connect     Connect two resources together"},
		{"create      Create Meroxa pipeline components"},
		{"describe    Describe a component"},
		{"help        Help about any command"},
		{"list        List components"},
		{"login       login or sign up to the Meroxa platform"},
		{"logout      logout of the Meroxa platform"},
		{"logs        Print logs for a component"},
		{"open        Open in a web browser"},
		{"remove      Remove a component"},
		{"update      Update a component"},
		{"version     Display the Meroxa CLI version"},
		{"Flags:\n" +
			"      --config string      config file (default is $HOME/meroxa.env)\n" +
			"      --debug              display any debugging information\n" +
			"  -h, --help               help for meroxa\n" +
			"      --json               output json\n" +
			"      --timeout duration   set the client timeout (default 10s)\n\n"},
		{"Use \"meroxa [command] --help\" for more information about a command."},
	}

	cmd := Cmd()
	var b bytes.Buffer
	cmd.SetOut(&b)
	cmd.Execute()
	out, err := ioutil.ReadAll(&b)

	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		if !strings.Contains(string(out), tt.expected) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.expected, string(out))
		}
	}
}
