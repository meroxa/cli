package apps

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/volatiletech/null/v8"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
	turbine "github.com/meroxa/turbine/init"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/mock"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/utils"
)

func TestDeployAppFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "path", required: false},
	}

	c := builder.BuildCobraCommand(&Deploy{})

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

const tempAppDir = "test-app"

func TestCreateApplication(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	name := "my-application"
	lang := GoLang

	// Create application locally
	path, err := initLocalApp(name)
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
	defer func() {
		_ = os.RemoveAll(tempAppDir)
	}()

	ai := &meroxa.CreateApplicationInput{
		Name:     name,
		Language: lang,
		GitSha:   "hardcoded",
		Pipeline: meroxa.EntityIdentifier{Name: null.StringFrom("default")},
	}

	a := &meroxa.Application{
		Name:     name,
		Language: lang,
		GitSha:   "hardcoded",
	}

	client.
		EXPECT().
		CreateApplication(
			ctx,
			ai,
		).
		Return(a, nil)

	d := &Deploy{
		client: client,
		logger: logger,
		path:   path,
		lang:   lang,
	}

	err = d.createApplication(ctx)

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

func initLocalApp(name string) (string, error) {
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
