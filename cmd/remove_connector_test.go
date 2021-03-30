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

func TestRemoveConnectorArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, errors.New("requires connector name\n\nUsage:\n  meroxa remove connector <name>"), ""},
		{[]string{"resName"}, nil, "resName"},
	}

	for _, tt := range tests {
		rc := RemoveConnector{}
		err := rc.setArgs(tt.args)

		if tt.err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != rc.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, rc.name)
		}
	}
}

func TestRemoveConnectorOutput(t *testing.T) {
	c := utils.GenerateConnector()

	output := utils.CaptureOutput(func() {
		rc := &RemoveConnector{}
		rc.output(&c)
	})

	expected := fmt.Sprintf("connector %s successfully removed", c.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestRemoveConnectorJSONOutput(t *testing.T) {
	c := utils.GenerateConnector()
	flagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		rc := &RemoveConnector{}
		rc.output(&c)
	})

	var parsedOutput meroxa.Connector
	json.Unmarshal([]byte(output), &parsedOutput)

	if !reflect.DeepEqual(c, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}
