package apps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestCreateApplicationArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires an application name"), name: ""},
		{args: []string{"application-name"}, err: nil, name: "application-name"},
	}

	for _, tt := range tests {
		cc := &Create{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestCreateApplicationFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "language", required: true},
	}

	c := builder.BuildCobraCommand(&Create{})

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		} else {
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
}

func TestCreateApplicationExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	name := "my-application"
	lang := NodeJs

	ai := &meroxa.CreateApplicationInput{
		Name:     name,
		Language: lang,
	}

	a := &meroxa.Application{
		Name:     name,
		Language: lang,
	}

	client.
		EXPECT().
		CreateApplication(
			ctx,
			ai,
		).
		Return(a, nil)

	c := &Create{
		client: client,
		logger: logger,
	}

	c.args.Name = ai.Name
	c.flags.Language = ai.Language

	err := c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating application %q...
Application %q successfully created!
`, name, name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotApplication meroxa.Application
	err = json.Unmarshal([]byte(gotJSONOutput), &gotApplication)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotApplication, *a) {
		t.Fatalf("expected \"%v\", got \"%v\"", *a, gotApplication)
	}
}
