package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
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
		hidden    bool
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
		{name: "env", required: false},
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

func TestCreateResourceExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
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
	c.flags.Type = string(r.Type)
	c.flags.URL = r.URL

	err := c.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating %q resource in %q environment...
Resource %q is successfully created!
`, cr.Type, meroxa.EnvironmentTypeCommon, cr.Name)

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

func TestCreateResourceExecutionWithEnvironmentName(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	c := &Create{
		client: client,
		logger: logger,
	}
	// Set up feature flags
	if global.Config == nil {
		build := builder.BuildCobraCommand(c)
		_ = global.PersistentPreRunE(build)
	}

	cr := utils.GenerateResourceWithEnvironment()

	r := meroxa.CreateResourceInput{
		Type:        "postgres",
		Name:        "",
		URL:         "https://foo.url",
		Credentials: nil,
		Metadata:    nil,
		Environment: &meroxa.EntityIdentifier{
			Name: cr.Environment.Name,
		},
	}

	client.
		EXPECT().
		CreateResource(
			ctx,
			&r,
		).
		Return(&cr, nil)

	c.args.Name = r.Name
	c.flags.Type = string(r.Type)
	c.flags.URL = r.URL
	c.flags.Environment = r.Environment.Name.String

	// override feature flags
	featureFlags := global.Config.Get(global.UserFeatureFlagsEnv)
	startingFlags := ""
	if featureFlags != nil {
		startingFlags = featureFlags.(string)
	}
	newFlags := ""
	if startingFlags != "" {
		newFlags = startingFlags + " "
	}
	newFlags += "environments"
	global.Config.Set(global.UserFeatureFlagsEnv, newFlags)

	err := c.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating %q resource in %q environment...
Resource %q is successfully created!
`, cr.Type, cr.Environment.Name.String, cr.Name)

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

	// Clear environments feature flags
	global.Config.Set(global.UserFeatureFlagsEnv, startingFlags)
}

func TestCreateResourceExecutionWithEnvironmentUUID(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	c := &Create{
		client: client,
		logger: logger,
	}
	// Set up feature flags
	if global.Config == nil {
		build := builder.BuildCobraCommand(c)
		_ = global.PersistentPreRunE(build)
	}

	cr := utils.GenerateResourceWithEnvironment()

	r := meroxa.CreateResourceInput{
		Type:        "postgres",
		Name:        "",
		URL:         "https://foo.url",
		Credentials: nil,
		Metadata:    nil,
		Environment: &meroxa.EntityIdentifier{
			UUID: cr.Environment.UUID,
		},
	}

	client.
		EXPECT().
		CreateResource(
			ctx,
			&r,
		).
		Return(&cr, nil)

	c.args.Name = r.Name
	c.flags.Type = string(r.Type)
	c.flags.URL = r.URL
	c.flags.Environment = r.Environment.UUID.String

	// override feature flags
	featureFlags := global.Config.Get(global.UserFeatureFlagsEnv)
	startingFlags := ""
	if featureFlags != nil {
		startingFlags = featureFlags.(string)
	}
	newFlags := ""
	if startingFlags != "" {
		newFlags = startingFlags + " "
	}
	newFlags += "environments"
	global.Config.Set(global.UserFeatureFlagsEnv, newFlags)

	err := c.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating %q resource in %q environment...
Resource %q is successfully created!
`, cr.Type, cr.Environment.UUID.String, cr.Name)

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

	// Clear environments feature flags
	global.Config.Set(global.UserFeatureFlagsEnv, startingFlags)
}

func TestCreateResourceExecutionWithEnvironmentUUIDWithoutFeatureFlag(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	c := &Create{
		client: client,
		logger: logger,
	}

	cr := utils.GenerateResourceWithEnvironment()

	r := meroxa.CreateResourceInput{
		Type:        "postgres",
		Name:        "",
		URL:         "https://foo.url",
		Credentials: nil,
		Metadata:    nil,
		Environment: &meroxa.EntityIdentifier{
			UUID: cr.Environment.UUID,
		},
	}

	c.args.Name = r.Name
	c.flags.Type = string(r.Type)
	c.flags.URL = r.URL
	c.flags.Environment = r.Environment.UUID.String

	err := c.Execute(ctx)
	if err == nil {
		t.Fatalf("unexpected success")
	}

	gotError := err.Error()
	wantError := `no access to the Meroxa self-hosted environments feature.
Sign up for the Beta here: https://share.hsforms.com/1Uq6UYoL8Q6eV5QzSiyIQkAc2sme`

	if gotError != wantError {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantError, gotError)
	}
}
