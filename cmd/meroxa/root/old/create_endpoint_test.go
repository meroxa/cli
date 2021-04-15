package old

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"github.com/spf13/cobra"
)

func TestCreateEndpointArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, nil, ""},
		{[]string{"name"}, nil, "name"},
	}

	for _, tt := range tests {
		ce := &CreateEndpoint{}
		err := ce.setArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != ce.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ce.name)
		}
	}
}

func TestCreateEndpointFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{"protocol", true, "p"},
		{"stream", true, "s"},
	}

	c := &cobra.Command{}
	ce := &CreateEndpoint{}
	ce.setFlags(c)

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

func TestCreateEndpointExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockCreateEndpointClient(ctrl)
	client.
		EXPECT().
		CreateEndpoint(
			ctx,
			"",
			"",
			"",
		).
		Return(nil)

	output := utils.CaptureOutput(func() {
		ce := &CreateEndpoint{}
		err := ce.execute(ctx, client)

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}
	})

	expected := fmt.Sprintf("Creating endpoint...")

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestCreateEndpointOutput(t *testing.T) {
	FlagRootOutputJSON = false

	output := utils.CaptureOutput(func() {
		ce := &CreateEndpoint{}
		ce.output()
	})

	expected := fmt.Sprintf("Endpoint successfully created!")

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}
