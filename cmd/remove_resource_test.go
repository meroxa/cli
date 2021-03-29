package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"reflect"
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

func TestRemoveResourceOutput(t *testing.T) {
	r := utils.GenerateResource()

	output := utils.CaptureOutput(func() {
		rr := &RemoveResource{}
		rr.output(&r)
	})

	expected := fmt.Sprintf("resource %s successfully removed", r.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestRemoveResourceJSONOutput(t *testing.T) {
	r := utils.GenerateResource()
	flagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		rr := &RemoveResource{}
		rr.output(&r)
	})

	var parsedOutput meroxa.Resource
	json.Unmarshal([]byte(output), &parsedOutput)

	if !reflect.DeepEqual(r, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}
