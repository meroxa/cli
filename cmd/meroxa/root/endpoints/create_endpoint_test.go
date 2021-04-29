package endpoints

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/utils"
)

func TestCreateEndpointArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: nil, name: ""},
		{args: []string{"endpoint-name"}, err: nil, name: "endpoint-name"},
	}

	for _, tt := range tests {
		cc := &CreateEndpoint{}
		err := cc.ParseArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestCreateEndpointFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "protocol", required: true, shorthand: "p", hidden: false},
		{name: "stream", required: true, shorthand: "s", hidden: false},
	}

	c := builder.BuildCobraCommand(&CreateEndpoint{})

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

		if cf.Hidden != f.hidden {
			if cf.Hidden {
				t.Fatalf("expected flag \"%s\" not to be hidden", f.name)
			} else {
				t.Fatalf("expected flag \"%s\" to be hidden", f.name)
			}
		}
	}
}

func TestCreateEndpointExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockCreateEndpointClient(ctrl)
	logger := log.NewTestLogger()

	client.
		EXPECT().
		CreateEndpoint(
			ctx,
			"",
			"",
			"",
		).
		Return(nil)

	ce := &CreateEndpoint{
		client: client,
		logger: logger,
	}

	err := ce.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprint(`Creating endpoint...
Endpoint successfully created!
`)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}
}
