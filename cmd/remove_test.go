package cmd

import (
	"bytes"
	"github.com/meroxa/cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io/ioutil"
	"strings"
	"testing"
)

var tests = []struct {
	expected string
}{
	{"Deprovision a component of the Meroxa platform, including pipelines,\n" +
		" resources, and connectors"},
	{"Usage:\n" +
		"  meroxa remove [command]"},
	{"Aliases:\n" +
		"  remove, rm, delete"},
	{"Available Commands:"},
	{"connector   Remove connector"},
	{"endpoint    Remove endpoint"},
	{"pipeline    Remove pipeline"},
	{"resource    Remove resource"},
	{"Flags:\n" +
		"  -h, --help   help for remove"},
}

func TestRemoveCmd(t *testing.T) {
	rootCmd := RootCmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"remove"})
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

func TestRemoveCmdWithRmAlias(t *testing.T) {
	rootCmd := RootCmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"rm"})
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

func TestRemoveCmdWithDeleteAlias(t *testing.T) {
	rootCmd := RootCmd()
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"delete"})
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

func TestRemoveFlags(t *testing.T) {
	expectedFlags := []struct {
		name       string
		required   bool
		shorthand  string
		persistent bool
	}{
		{"force", false, "f", true},
	}

	c := &cobra.Command{}
	r := &Remove{}
	r.setFlags(c)

	for _, f := range expectedFlags {
		var cf *pflag.Flag

		if f.persistent {
			cf = c.PersistentFlags().Lookup(f.name)
		} else {
			cf = c.Flags().Lookup(f.name)
		}

		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		}

		if f.shorthand != cf.Shorthand {
			t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
		}

		if f.required && !utils.IsFlagRequired(cf) {
			t.Fatalf("expected flag \"%s\" to be required", f.name)
		}
	}
}
