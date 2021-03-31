package cmd

import (
	"errors"
	"github.com/meroxa/cli/utils"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func TestUpdatePipelineArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, errors.New("requires pipeline name\n\nUsage:\n  meroxa update pipeline <name> --state <pause|resume|restart>"), ""},
		{[]string{"pipelineName"}, nil, "pipelineName"},
	}

	for _, tt := range tests {
		up := &UpdatePipeline{}
		err := up.setArgs(tt.args)

		if err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err.Error(), err.Error())
		}

		if tt.name != up.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, up.name)
		}
	}
}

func TestUpdatePipelineFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{"state", true, ""},
	}

	c := &cobra.Command{}
	up := &UpdatePipeline{}
	up.setFlags(c)

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
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
