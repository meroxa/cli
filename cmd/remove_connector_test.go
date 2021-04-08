package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
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
		{nil, errors.New("requires connector name\n\nUsage:\n  meroxa remove connector NAME"), ""},
		{[]string{"resName"}, nil, "resName"},
	}

	r := &Remove{}
	for _, tt := range tests {
		rc := &RemoveConnector{removeCmd: r}
		err := rc.setArgs(tt.args)

		if tt.err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != rc.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, rc.name)
		}

		if err == nil {
			componentType := "connector"
			if rc.removeCmd.componentType != componentType {
				t.Fatalf("expected type to be set to %q", componentType)
			}

			if rc.removeCmd.confirmableName != rc.name {
				t.Fatalf("expected \"confirmableName\" to be set to %q", rc.name)
			}
		}
	}
}

func TestRemoveConnectorExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockRemoveConnectorClient(ctrl)

	c := utils.GenerateConnector(0)

	client.
		EXPECT().
		GetConnectorByName(ctx, c.Name).
		Return(&c, nil).
		MaxTimes(2)

	client.
		EXPECT().
		DeleteConnector(ctx, c.ID).
		Return(nil).
		MaxTimes(2)

	r := &Remove{}

	rc := &RemoveConnector{
		name:      c.Name,
		removeCmd: r,
	}
	got, err := rc.execute(ctx, client)

	if !reflect.DeepEqual(got, &c) {
		t.Fatalf("expected \"%v\", got \"%v\"", &c, got)
	}

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestRemoveConnectorOutput(t *testing.T) {
	c := utils.GenerateConnector(0)
	flagRootOutputJSON = false

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
	c := utils.GenerateConnector(0)
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
