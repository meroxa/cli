package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		{name: "url", shorthand: "u"},
		{name: "username", required: false, shorthand: ""},
		{name: "password", required: false, shorthand: ""},
		{name: "ca-cert", required: false, shorthand: ""},
		{name: "client-cert", required: false, shorthand: ""},
		{name: "client-key", required: false, shorthand: ""},
		{name: "ssl", required: false, shorthand: ""},
		{name: "metadata", required: false, shorthand: "m"},
		{name: "env", required: false},
		{name: "token", required: false},
		{name: "ssh-url", required: false},
		{name: "ssh-private-key", required: false},
		{name: "private-key-file", required: false},
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
	c.flags.Environment = r.Environment.Name

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
`, cr.Type, cr.Environment.Name, cr.Name)

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
	c.flags.Environment = r.Environment.UUID

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
`, cr.Type, cr.Environment.UUID, cr.Name)

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

	if global.Config == nil {
		build := builder.BuildCobraCommand(c)
		_ = global.PersistentPreRunE(build)
	}
	global.Config.Set(global.UserFeatureFlagsEnv, "")

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
	c.flags.Environment = r.Environment.UUID

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

func TestCreateResourceExecutionPrivateKeyFlags(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()

	keyVal := "super-secret"
	keyFile := filepath.Join("/tmp", uuid.NewString())
	err := os.WriteFile(keyFile, []byte(keyVal), 0o600)
	require.NoError(t, err)

	tests := []struct {
		name                     string
		inputType                string
		inputSSHPrivateKeyFlag   string
		inputPasswordFlag        string
		inputPrivateKeyFileFlag  string
		expectedPassword         string
		expectedSSHPrivateKeyVal string
	}{
		{
			name:                     "create postgres resource with SSH Tunnel --ssh-private-key",
			inputType:                string(meroxa.ResourceTypePostgres),
			inputSSHPrivateKeyFlag:   keyVal,
			expectedPassword:         "",
			expectedSSHPrivateKeyVal: keyVal,
		},
		{
			name:                     "create postgres resource with SSH Tunnel --private-key-file",
			inputType:                string(meroxa.ResourceTypePostgres),
			inputPrivateKeyFileFlag:  keyFile,
			expectedPassword:         "",
			expectedSSHPrivateKeyVal: keyVal,
		},
		{
			name:                     "create postgres resource with both SSH flags",
			inputType:                string(meroxa.ResourceTypePostgres),
			inputPrivateKeyFileFlag:  keyFile,
			inputSSHPrivateKeyFlag:   keyVal,
			expectedPassword:         "",
			expectedSSHPrivateKeyVal: keyVal,
		},
		{
			name:                     "create snowflake resource with --password",
			inputType:                string(meroxa.ResourceTypeSnowflake),
			inputPasswordFlag:        keyVal,
			expectedPassword:         keyVal,
			expectedSSHPrivateKeyVal: "",
		},
		{
			name:                     "create snowflake resource with --private-key-file",
			inputPrivateKeyFileFlag:  keyFile,
			inputType:                string(meroxa.ResourceTypeSnowflake),
			expectedPassword:         keyVal,
			expectedSSHPrivateKeyVal: keyVal,
		},
		{
			name:                     "create snowflake resource with both secret flags",
			inputPasswordFlag:        keyVal,
			inputPrivateKeyFileFlag:  keyFile,
			inputType:                string(meroxa.ResourceTypeSnowflake),
			expectedPassword:         keyVal,
			expectedSSHPrivateKeyVal: keyVal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := mock.NewMockClient(ctrl)

			c := &Create{
				client: client,
				logger: logger,
			}

			client.
				EXPECT().
				CreateResource(
					ctx,
					gomock.Any(),
				).
				Return(&meroxa.Resource{}, nil)

			c.args.Name = "my-resource"
			c.flags.Type = tc.inputType
			c.flags.URL = "anything"
			c.flags.Password = tc.inputPasswordFlag
			c.flags.SSHPrivateKey = tc.inputSSHPrivateKeyFlag
			c.flags.PrivateKeyFile = tc.inputPrivateKeyFileFlag

			err := c.Execute(ctx)
			if err != nil {
				t.Fatalf("not expected error, got %q", err.Error())
			}

			assert.Equalf(t, tc.expectedSSHPrivateKeyVal, c.flags.SSHPrivateKey, "mistach in private key flag value")
			assert.Equalf(t, tc.expectedPassword, c.flags.Password, "mismatch in password flag value")
		})
	}
}

//nolint:funlen
func TestCreateResourceURLFlag(t *testing.T) {
	resourceName := "my-resource"
	tests := []struct {
		description  string
		resourceType string
		url          string
		client       func(*gomock.Controller) *mock.MockClient
		wantOutput   string
		wantErr      error
	}{
		{
			description:  "Do not require URL for Notion",
			resourceType: string(meroxa.ResourceTypeNotion),
			client: func(ctrl *gomock.Controller) *mock.MockClient {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().CreateResource(gomock.Any(), gomock.Any()).Return(&meroxa.Resource{Name: resourceName}, nil).Times(1)
				return client
			},
			wantOutput: `Creating "notion" resource in "common" environment...
Resource "my-resource" is successfully created!
`,
		},
		{
			description:  "Allow default URL value for for Notion",
			resourceType: string(meroxa.ResourceTypeNotion),
			url:          "https://api.notion.com",
			client: func(ctrl *gomock.Controller) *mock.MockClient {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().CreateResource(gomock.Any(), gomock.Any()).Return(&meroxa.Resource{Name: resourceName}, nil).Times(1)
				return client
			},
			wantOutput: `Creating "notion" resource in "common" environment...
Resource "my-resource" is successfully created!
`,
		},
		{
			description:  "Warn about non-default URL value for for Notion",
			resourceType: string(meroxa.ResourceTypeNotion),
			url:          "https://wild.west.api.notion.com",
			client: func(ctrl *gomock.Controller) *mock.MockClient {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().CreateResource(gomock.Any(), gomock.Any()).Return(&meroxa.Resource{Name: resourceName}, nil).Times(1)
				return client
			},
			wantOutput: `Ignoring API URL override (https://wild.west.api.notion.com) for Notion resource configuration.
Creating "notion" resource in "common" environment...
Resource "my-resource" is successfully created!
`,
		},
		{
			description:  "Do not require URL for Spire Maritime AIS",
			resourceType: string(meroxa.ResourceTypeSpireMaritimeAIS),
			client: func(ctrl *gomock.Controller) *mock.MockClient {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().CreateResource(gomock.Any(), gomock.Any()).Return(&meroxa.Resource{Name: resourceName}, nil).Times(1)
				return client
			},
			wantOutput: `Creating "spire_maritime_ais" resource in "common" environment...
Resource "my-resource" is successfully created!
`,
		},
		{
			description:  "Allow default URL for Spire Maritime AIS",
			resourceType: string(meroxa.ResourceTypeSpireMaritimeAIS),
			url:          "https://api.spire.com/graphql",
			client: func(ctrl *gomock.Controller) *mock.MockClient {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().CreateResource(gomock.Any(), gomock.Any()).Return(&meroxa.Resource{Name: resourceName}, nil).Times(1)
				return client
			},
			wantOutput: `Creating "spire_maritime_ais" resource in "common" environment...
Resource "my-resource" is successfully created!
`,
		},
		{
			description:  "Warn about non-default URL for Spire Maritime AIS",
			resourceType: string(meroxa.ResourceTypeSpireMaritimeAIS),
			url:          "https://api.spire.com/ascii",
			client: func(ctrl *gomock.Controller) *mock.MockClient {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().CreateResource(gomock.Any(), gomock.Any()).Return(&meroxa.Resource{Name: resourceName}, nil).Times(1)
				return client
			},
			wantOutput: `Ignoring API URL override (https://api.spire.com/ascii) for Spire Maritime AIS resource configuration.
Creating "spire_maritime_ais" resource in "common" environment...
Resource "my-resource" is successfully created!
`,
		},
		{
			description:  "Require URL for one of the rest of the types",
			resourceType: string(meroxa.ResourceTypePostgres),
			client: func(ctrl *gomock.Controller) *mock.MockClient {
				client := mock.NewMockClient(ctrl)
				return client
			},
			wantErr:    fmt.Errorf(`required flag(s) "url" not set`),
			wantOutput: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			logger := log.NewTestLogger()
			c := &Create{
				client: tc.client(ctrl),
				logger: logger,
			}
			c.args.Name = resourceName
			c.flags.Type = tc.resourceType
			c.flags.URL = tc.url

			err := c.Execute(ctx)
			gotLeveledOutput := logger.LeveledOutput()
			assert.Equal(t, tc.wantOutput, gotLeveledOutput)

			if err != nil {
				if tc.wantErr == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
