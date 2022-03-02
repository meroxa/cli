package apps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
	turbine "github.com/meroxa/turbine/init"
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
		{name: "path", required: false},
		{name: "lang", required: false, shorthand: "l"},
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

func TestCreateApplicationExecutionNoLangNoPath(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	name := "my-application"

	c := &Create{
		client: client,
		logger: logger,
	}

	c.args.Name = name
	err := c.Execute(ctx)

	expectedError := "language is required either using --path ~/turbine/my-app or --lang. Type `meroxa help apps create` for more information"
	if err != nil && err.Error() != expectedError {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

const tempAppDir = "test-app"

func TestCreateApplicationExecutionWithAppJSON(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	name := "my-application"
	lang := GoLang

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

	path, err := initApp(name)
	defer func() {
		_ = os.RemoveAll(tempAppDir)
	}()
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	c.args.Name = ai.Name
	c.flags.Path = path

	err = c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating application %q with language %q...
Application %q successfully created!
`, name, lang, name)

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

func initApp(name string) (string, error) {
	if err := os.Mkdir(tempAppDir, 0700); err != nil {
		return "", err
	}

	if err := turbine.Init(name, tempAppDir); err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", cwd, tempAppDir, name), nil
}

func TestCreateApplicationExecutionWithLang(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	name := "my-application"
	lang := GoLang

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
	c.flags.Lang = lang

	err := c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating application %q with language %q...
Application %q successfully created!
`, name, lang, name)

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
