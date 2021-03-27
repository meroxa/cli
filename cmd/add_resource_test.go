package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	utils "github.com/meroxa/cli/utils"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
	"reflect"
	"strings"
	"testing"
)

func TestAddResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err error
		name string
	}{
		{[]string{""},nil, ""},
		{[]string{"resName"},nil, "resName"},
	}

	for _, tt := range tests {
		name, err := AddResource{}.checkArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, name)
		}
	}
}

func TestAddResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name string
		required bool
		shorthand string
	}{
		{"type", true, ""},
		{"url", true, "u"},
		{"credentials", false, ""},
		{"metadata", false, "m"},
	}

	c := &cobra.Command{}
	AddResource{}.setFlags(c)

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

func TestAddResourceOutput(t *testing.T) {
	r := utils.GenerateResource()

	output := utils.CaptureOutput(func() {
		AddResource{}.output(&r)
	})

	expected := fmt.Sprintf("Resource %s successfully added!", r.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestAddResourceJSONOutput(t *testing.T) {
	r := utils.GenerateResource()
	flagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		AddResource{}.output(&r)
	})

	var parsedOutput meroxa.Resource
	json.Unmarshal([]byte(output), &parsedOutput)


	if !reflect.DeepEqual(r, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}

func TestAddResourceExecution(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	client := mock.NewMockAddResourceClient(ctrl)

	r := meroxa.CreateResourceInput{
		Type:        "postgres",
		Name:        "",
		URL:         "https://foo.url",
		Credentials: nil,
		Metadata:    nil,
	}

	client.
		EXPECT().
		CreateResource(
			ctx,
			gomock.Eq(&r),
		).
		DoAndReturn(func() (*meroxa.Resource, error) {
			nr := utils.GenerateResource()
			return &nr, nil
		})

	got, err := AddResource{}.execute(ctx, client, r)

	if got != nil {
		t.Fatal("not good")
	}
	if err == nil {
		t.Fatal("not good")
	}
}
