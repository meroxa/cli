package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"testing"
)

func TestAddResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err error
		name string
	}{
		{[]string{""},nil, ""},
		{[]string{"resName"},nil, "resName"},
	}

	for _, tt := range tests {
		name, err := addResourceArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, name)
		}
	}
}

func isFlagRequired(flag *pflag.Flag) bool{
	requiredAnnotation := "cobra_annotation_bash_completion_one_required_flag"

	if len(flag.Annotations[requiredAnnotation]) > 0 && flag.Annotations[requiredAnnotation][0] == "true" {
		return true
	}

	return false
}

func TestAddResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name string
		required bool
		shorthand string
	}{
		{"type", true, ""},
		{"url", true, "u"},
		{"credentials", false, ""},
		{"metadata", false, "m"},
	}

	c := addResourceFlags(&cobra.Command{})

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		}

		if f.shorthand != cf.Shorthand {
			t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
		}

		if f.required && !isFlagRequired(cf) {
			t.Fatalf("expected flag \"%s\" to be required", f.name)
		}
	}
}

// TODO: Test adddResource

// TODO Test printOutResource
// given a resource Type, you get "successfully added" when not --json


