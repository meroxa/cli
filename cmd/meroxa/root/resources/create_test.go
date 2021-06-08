package resources

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

func TestCreateResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, nil, ""},
		{[]string{"my-resource"}, nil, "my-resource"},
	}

	for _, tt := range tests {
		c := &Create{}
		err := c.ParseArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != c.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, c.args.Name)
		}
	}
}

func TestCreateResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{name: "type", required: true, shorthand: ""},
		{name: "url", required: true, shorthand: "u"},
		{name: "username", required: false, shorthand: ""},
		{name: "password", required: false, shorthand: ""},
		{name: "ca-cert", required: false, shorthand: ""},
		{name: "client-cert", required: false, shorthand: ""},
		{name: "client-key", required: false, shorthand: ""},
		{name: "ssl", required: false, shorthand: ""},
		{name: "metadata", required: false, shorthand: "m"},
	}

	c := builder.BuildCobraCommand(&Create{})

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

func TestCreateResourceExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockCreateResourceClient(ctrl)
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

	c := &Create{
		client: client,
		logger: logger,
	}
	c.args.Name = r.Name
	c.flags.Type = r.Type
	c.flags.URL = r.URL

	err := c.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating %q resource...
%q resource is successfully created!
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
