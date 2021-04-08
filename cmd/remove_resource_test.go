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

func TestRemoveResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, errors.New("requires resource name\n\nUsage:\n  meroxa remove resource NAME"), ""},
		{[]string{"resName"}, nil, "resName"},
	}

	r := &Remove{}
	for _, tt := range tests {
		rr := &RemoveResource{removeCmd: r}
		err := rr.setArgs(tt.args)

		if tt.err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != rr.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, rr.name)
		}

		if err == nil {
			componentType := "resource"
			if rr.removeCmd.componentType != componentType {
				t.Fatalf("expected type to be set to %q", componentType)
			}

			if rr.removeCmd.confirmableName != rr.name {
				t.Fatalf("expected \"confirmableName\" to be set to %q", rr.name)
			}
		}
	}
}

func TestRemoveResourceExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockRemoveResourceClient(ctrl)

	r := utils.GenerateResource()

	client.
		EXPECT().
		GetResourceByName(ctx, r.Name).
		Return(&r, nil)

	client.
		EXPECT().
		DeleteResource(ctx, r.ID).
		Return(nil)

	rc := &Remove{}

	rr := &RemoveResource{
		name:      r.Name,
		removeCmd: rc,
	}
	got, err := rr.execute(ctx, client)

	if !reflect.DeepEqual(got, &r) {
		t.Fatalf("expected \"%v\", got \"%v\"", &r, got)
	}

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestRemoveResourceOutput(t *testing.T) {
	r := utils.GenerateResource()
	flagRootOutputJSON = false

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
