package resources

import (
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/utils"
)

func TestCreateResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, nil, ""},
		{[]string{"my-resource"}, nil, "my-resource"},
	}

	for _, tt := range tests {
		c := &Create{}
		err := c.ParseArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != c.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, c.args.Name)
		}
	}
}

func TestCreateResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "type", required: true, shorthand: ""},
		{name: "url", required: true, shorthand: "u"},
		{name: "username", required: false, shorthand: ""},
		{name: "password", required: false, shorthand: ""},
		{name: "ca-cert", required: false, shorthand: ""},
		{name: "client-cert", required: false, shorthand: ""},
		{name: "client-key", required: false, shorthand: ""},
		{name: "ssl", required: false, shorthand: ""},
		{name: "metadata", required: false, shorthand: "m"},
		{name: "env", required: false},
	}

	c := builder.BuildCobraCommand(&Create{})

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		} else {
			if f.shorthand != cf.Shorthand {
				t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
			}

			if f.required && !utils.IsFlagRequired(cf) {
				t.Fatalf("expected flag \"%s\" to be required", f.name)
			}

			if cf.Hidden != f.hidden {
				if cf.Hidden {
					t.Fatalf("expected flag \"%s\" not to be hidden", f.name)
				} else {
					t.Fatalf("expected flag \"%s\" to be hidden", f.name)
				}
			}
		}
	}
}
