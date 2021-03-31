package cmd

import (
	"errors"
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
