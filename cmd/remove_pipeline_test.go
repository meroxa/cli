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

func TestRemovePipelineArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, errors.New("requires pipeline name\n\nUsage:\n  meroxa remove pipeline NAME"), ""},
		{[]string{"endpoint-name"}, nil, "endpoint-name"},
	}

	r := &Remove{}
	for _, tt := range tests {
		rr := &RemovePipeline{removeCmd: r}
		err := rr.setArgs(tt.args)

		if tt.err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != rr.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, rr.name)
		}

		if err == nil {
			componentType := "pipeline"
			if rr.removeCmd.componentType != componentType {
				t.Fatalf("expected type to be set to %q", componentType)
			}

			if rr.removeCmd.confirmableName != rr.name {
				t.Fatalf("expected \"confirmableName\" to be set to %q", rr.name)
			}
		}
	}
}

func TestRemovePipelineExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockRemovePipelineClient(ctrl)

	p := utils.GeneratePipeline()

	client.
		EXPECT().
		GetPipelineByName(ctx, p.Name).
		Return(&p, nil)

	client.
		EXPECT().
		DeletePipeline(ctx, p.ID).
		Return(nil)

	rc := &Remove{}

	rr := &RemovePipeline{
		name:      p.Name,
		removeCmd: rc,
	}
	got, err := rr.execute(ctx, client)

	if !reflect.DeepEqual(got, &p) {
		t.Fatalf("expected \"%v\", got \"%v\"", &p, got)
	}

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestRemovePipelineOutput(t *testing.T) {
	p := utils.GeneratePipeline()
	flagRootOutputJSON = false

	output := utils.CaptureOutput(func() {
		rr := &RemovePipeline{}
		rr.output(&p)
	})

	expected := fmt.Sprintf("pipeline %s successfully removed", p.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestRemovePipelineJSONOutput(t *testing.T) {
	r := utils.GeneratePipeline()
	flagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		rr := &RemovePipeline{}
		rr.output(&r)
	})

	var parsedOutput meroxa.Pipeline
	json.Unmarshal([]byte(output), &parsedOutput)

	if !reflect.DeepEqual(r, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}
