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
	"github.com/spf13/cobra"
	"reflect"
	"strings"
	"testing"
)

func TestUpdatePipelineArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, errors.New("requires pipeline name"), ""},
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
		{"state", false, ""},
		{"name", false, ""},
		{"metadata", false, "m"},
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

func TestUpdatePipelineExecutionNoFlags(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdatePipelineClient(ctrl)

	up := &UpdatePipeline{}
	_, err := up.execute(ctx, client)

	expected := "requires either --name, --state or --metadata"

	if err != nil && err.Error() != expected {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestUpdatePipelineExecutionWithNewState(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdatePipelineClient(ctrl)

	flagRootOutputJSON = false

	p := utils.GeneratePipeline()
	newState := "pause"

	client.
		EXPECT().
		GetPipelineByName(ctx, p.Name).
		Return(&p, nil)

	client.
		EXPECT().
		UpdatePipelineStatus(ctx, p.ID, newState).
		Return(&p, nil)

	output := utils.CaptureOutput(func() {
		up := &UpdatePipeline{}
		up.name = p.Name
		up.state = newState

		got, err := up.execute(ctx, client)

		if !reflect.DeepEqual(got, &p) {
			t.Fatalf("expected \"%v\", got \"%v\"", &p, got)
		}

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}
	})

	expected := fmt.Sprintf("Updating %s pipeline...", p.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestUpdatePipelineExecutionWithNewName(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdatePipelineClient(ctrl)

	flagRootOutputJSON = false

	p := utils.GeneratePipeline()
	newName := "new-pipeline-name"
	pi := meroxa.UpdatePipelineInput{
		Name: newName,
	}

	client.
		EXPECT().
		GetPipelineByName(ctx, p.Name).
		Return(&p, nil)

	client.
		EXPECT().
		UpdatePipeline(ctx, p.ID, pi).
		Return(&p, nil)

	output := utils.CaptureOutput(func() {
		up := &UpdatePipeline{}
		up.name = p.Name

		// What we're trying to update
		up.newName = newName

		got, err := up.execute(ctx, client)

		if !reflect.DeepEqual(got, &p) {
			t.Fatalf("expected \"%v\", got \"%v\"", &p, got)
		}

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}
	})

	expected := fmt.Sprintf("Updating %s pipeline...", p.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestUpdatePipelineExecutionWithNewMetadata(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdatePipelineClient(ctrl)

	flagRootOutputJSON = false

	p := utils.GeneratePipeline()
	var pi meroxa.UpdatePipelineInput

	pi.Metadata = map[string]string{"key": "value"}

	client.
		EXPECT().
		GetPipelineByName(ctx, p.Name).
		Return(&p, nil)

	client.
		EXPECT().
		UpdatePipeline(ctx, p.ID, pi).
		Return(&p, nil)

	output := utils.CaptureOutput(func() {
		up := &UpdatePipeline{}
		up.name = p.Name

		// What we're trying to update
		up.metadata = "{\"key\": \"value\"}"

		got, err := up.execute(ctx, client)

		if !reflect.DeepEqual(got, &p) {
			t.Fatalf("expected \"%v\", got \"%v\"", &p, got)
		}

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}
	})

	expected := fmt.Sprintf("Updating %s pipeline...", p.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestUpdatePipelineOutput(t *testing.T) {
	p := utils.GeneratePipeline()
	flagRootOutputJSON = false

	output := utils.CaptureOutput(func() {
		up := &UpdatePipeline{}
		up.output(&p)
	})

	expected := fmt.Sprintf("pipeline %s successfully updated!", p.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestUpdatePipelineJSONOutput(t *testing.T) {
	r := utils.GeneratePipeline()
	flagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		ar := &UpdatePipeline{}
		ar.output(&r)
	})

	var parsedOutput meroxa.Pipeline
	json.Unmarshal([]byte(output), &parsedOutput)

	if !reflect.DeepEqual(r, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}
