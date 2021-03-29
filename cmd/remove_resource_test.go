package cmd

import (
	"errors"
	"strings"
	"testing"
)

func TestRemoveResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, errors.New("requires resource name\n\nUsage:\n  meroxa remove resource <name>"), ""},
		{[]string{"resName"}, nil, "resName"},
	}

	for _, tt := range tests {
		rr := RemoveResource{}
		err := rr.setArgs(tt.args)

		if tt.err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != rr.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, rr.name)
		}
	}
}
