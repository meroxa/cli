package add

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

func TestAddResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, nil, ""},
		{[]string{"resName"}, nil, "resName"},
	}

	for _, tt := range tests {
		ar := &AddResource{}
		err := ar.ParseArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != ar.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ar.args.Name)
		}
	}
}

func TestAddResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{"type", true, ""},
		{"url", true, "u"},
		{"username", false, ""},
		{"password", false, ""},
		{"ca-cert", false, ""},
		{"client-cert", false, ""},
		{"client-key", false, ""},
		{"ssl", false, ""},
		{"metadata", false, "m"},
	}

	c := builder.BuildCobraCommand(&AddResource{})

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

func TestAddResourceExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockAddResourceClient(ctrl)
	logger := log.NewTestLogger()

	r := meroxa.CreateResourceInput{
		Type:        "postgres",
		Name:        "",
		URL:         "https://foo.url",
		Credentials: nil,
		Metadata:    nil,
	}

	cr := utils.GenerateResource()
	client.
		EXPECT().
		CreateResource(
			ctx,
			&r,
		).
		Return(&cr, nil)

	ar := &AddResource{
		client: client,
		logger: logger,
	}
	ar.args.Name = r.Name
	ar.flags.Type = r.Type
	ar.flags.Url = r.URL

	err := ar.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Adding postgres resource...
%s resource with name %s successfully added!
`, cr.Type, cr.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotResource meroxa.Resource
	err = json.Unmarshal([]byte(gotJSONOutput), &gotResource)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotResource, cr) {
		t.Fatalf("expected \"%v\", got \"%v\"", cr, gotResource)
	}
}
