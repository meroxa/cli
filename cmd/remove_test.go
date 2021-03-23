package cmd

import (
	"bytes"
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
