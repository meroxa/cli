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

func TestUpdateResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, errors.New("requires resource name"), ""},
		{[]string{"resName"}, nil, "resName"},
	}

	for _, tt := range tests {
		ur := &UpdateResource{}
		err := ur.setArgs(tt.args)

		if err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err.Error(), err.Error())
		}

		if tt.name != ur.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ur.name)
		}
	}
}

func TestUpdateResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{"credentials", false, ""},
		{"metadata", false, "m"},
		{"name", false, ""},
		{"url", false, "u"},
	}

	c := &cobra.Command{}
	ur := &UpdateResource{}
	ur.setFlags(c)

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

func TestUpdateResourceExecutionWithNewName(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateResourceClient(ctrl)

	flagRootOutputJSON = false

	r := utils.GenerateResource()

	newName := "my-new-resource-name"
	nr := meroxa.UpdateResourceInput{
		Name: newName,
	}

	client.
		EXPECT().
		UpdateResource(ctx, r.Name, nr).
		Return(&r, nil)

	output := utils.CaptureOutput(func() {
		ur := &UpdateResource{
			name:    r.Name,
			newName: newName,
		}

		got, err := ur.execute(ctx, client)

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}

		if !reflect.DeepEqual(got, &r) {
			t.Fatalf("expected \"%v\", got \"%v\"", &r, got)
		}

	})

	expected := fmt.Sprintf("Updating %s resource...", r.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestUpdateResourceExecutionWithNewMetadata(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateResourceClient(ctrl)

	flagRootOutputJSON = false

	r := utils.GenerateResource()

	ur := &UpdateResource{
		name:     r.Name,
		metadata: `{"metakey":"metavalue"}`,
	}

	var metadata map[string]string

	json.Unmarshal([]byte(ur.metadata), &metadata)
	nr := meroxa.UpdateResourceInput{
		Metadata: metadata,
	}

	client.
		EXPECT().
		UpdateResource(ctx, ur.name, nr).
		Return(&r, nil)

	output := utils.CaptureOutput(func() {
		got, err := ur.execute(ctx, client)

		if !reflect.DeepEqual(got, &r) {
			t.Fatalf("expected \"%v\", got \"%v\"", &r, got)
		}

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}
	})

	expected := fmt.Sprintf("Updating %s resource...", r.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestUpdateResourceExecutionWithNewUrl(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateResourceClient(ctrl)

	flagRootOutputJSON = false

	r := utils.GenerateResource()

	ur := &UpdateResource{
		name: r.Name,
		url:  "https://newUrl.io",
	}

	nr := meroxa.UpdateResourceInput{
		URL: ur.url,
	}

	client.
		EXPECT().
		UpdateResource(ctx, ur.name, nr).
		Return(&r, nil)

	output := utils.CaptureOutput(func() {
		got, err := ur.execute(ctx, client)

		if !reflect.DeepEqual(got, &r) {
			t.Fatalf("expected \"%v\", got \"%v\"", &r, got)
		}

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}
	})

	expected := fmt.Sprintf("Updating %s resource...", r.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestUpdateResourceExecutionWithNewCredentials(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateResourceClient(ctrl)

	flagRootOutputJSON = false

	r := utils.GenerateResource()

	newCred := meroxa.Credentials{Username: "newUsername"}
	nr := meroxa.UpdateResourceInput{
		Credentials: &newCred,
	}

	client.
		EXPECT().
		UpdateResource(ctx, r.Name, nr).
		Return(&r, nil)

	output := utils.CaptureOutput(func() {
		ur := &UpdateResource{
			name:        r.Name,
			credentials: "{\"username\":\"newUsername\"}",
		}

		got, err := ur.execute(ctx, client)

		if !reflect.DeepEqual(got, &r) {
			t.Fatalf("expected \"%v\", got \"%v\"", &r, got)
		}

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}
	})

	expected := fmt.Sprintf("Updating %s resource...", r.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestUpdateResourceExecutionNoFlags(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateResourceClient(ctrl)

	ur := &UpdateResource{}
	_, err := ur.execute(ctx, client)

	expected := "requires either `--credentials`, `--name`, `--metadata` or `--url` to update the resource"

	if err != nil && err.Error() != expected {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestUpdateResourceOutput(t *testing.T) {
	r := utils.GenerateResource()
	flagRootOutputJSON = false

	output := utils.CaptureOutput(func() {
		ur := &UpdateResource{name: r.Name}
		ur.output(&r)
	})

	expected := fmt.Sprintf("Resource %s successfully updated!", r.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestUpdateResourceJSONOutput(t *testing.T) {
	r := utils.GenerateResource()
	flagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		ur := &UpdateResource{}
		ur.output(&r)
	})

	var parsedOutput meroxa.Resource
	json.Unmarshal([]byte(output), &parsedOutput)

	if !reflect.DeepEqual(r, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}
