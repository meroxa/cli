package apps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/meroxa/cli/config"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
	turbine "github.com/meroxa/turbine-go/init"

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
		{name: "docker-hub-username", required: false, hidden: true},
		{name: "docker-hub-access-token", required: false, hidden: true},
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

func TestValidateDockerHubFlags(t *testing.T) {
	tests := []struct {
		desc                 string
		dockerHubUserName    string
		dockerHubAccessToken string
		err                  error
	}{
		{
			desc:                 "Neither DockerHub flags are present",
			dockerHubUserName:    "",
			dockerHubAccessToken: "",
			err:                  nil,
		},
		{
			desc:                 "DockerHubUserName is specified, but DockerHubAccessToken isn't",
			dockerHubUserName:    "my-docker-hub-username",
			dockerHubAccessToken: "",
			err:                  errors.New("--docker-hub-access-token is required when --docker-hub-username is present"),
		},
		{
			desc:                 "DockerHubAccessToken is specified, but DockerHubUserName isn't",
			dockerHubUserName:    "",
			dockerHubAccessToken: "my-docker-hub-access-token",
			err:                  errors.New("--docker-hub-username is required when --docker-hub-access-token is present"),
		},
		{
			desc:                 "BothDockerHubAccessToken and DockerHubUserName are specified",
			dockerHubUserName:    "my-docker-hub-username",
			dockerHubAccessToken: "my-docker-hub-access-token",
			err:                  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			d := &Deploy{}
			d.flags.DockerHubUserName = tc.dockerHubUserName
			d.flags.DockerHubAccessToken = tc.dockerHubAccessToken
			err := d.validateDockerHubFlags()

			if err != nil && tc.err.Error() != err.Error() {
				t.Fatalf("expected %v, got %v", tc.err, err)
			}

			if err == nil {
				if d.localDeploy.DockerHubUserNameEnv != tc.dockerHubUserName {
					t.Fatalf("expected DockerHubUserNameEnv to be %q, got %q", tc.dockerHubUserName, d.localDeploy.DockerHubUserNameEnv)
				}

				if d.localDeploy.DockerHubAccessTokenEnv != tc.dockerHubAccessToken {
					t.Fatalf("expected DockerHubAccessTokenEnv to be %q, got %q", tc.dockerHubAccessToken, d.localDeploy.DockerHubAccessTokenEnv)
				}
			}
		})
	}
}

func TestValidateDockerHubEnVars(t *testing.T) {
	tests := []struct {
		desc                 string
		dockerHubUserName    string
		dockerHubAccessToken string
		err                  error
	}{
		{
			desc:                 "Neither DockerHub flags are present",
			dockerHubUserName:    "",
			dockerHubAccessToken: "",
			err:                  nil,
		},
		{
			desc:                 "DockerHubUserName is specified, but DockerHubAccessToken isn't",
			dockerHubUserName:    "my-docker-hub-username",
			dockerHubAccessToken: "",
			err:                  fmt.Errorf("%s is required when %s is present", dockerHubAccessTokenEnv, dockerHubUserNameEnv),
		},
		{
			desc:                 "DockerHubAccessToken is specified, but DockerHubUserName isn't",
			dockerHubUserName:    "",
			dockerHubAccessToken: "my-docker-hub-access-token",
			err:                  fmt.Errorf("%s is required when %s is present", dockerHubUserNameEnv, dockerHubAccessTokenEnv),
		},
		{
			desc:                 "BothDockerHubAccessToken and DockerHubUserName are specified",
			dockerHubUserName:    "my-docker-hub-username",
			dockerHubAccessToken: "my-docker-hub-access-token",
			err:                  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := config.NewInMemoryConfig()

			d := &Deploy{
				config: cfg,
			}
			d.config.Set(dockerHubUserNameEnv, tc.dockerHubUserName)
			d.config.Set(dockerHubAccessTokenEnv, tc.dockerHubAccessToken)
			err := d.validateDockerHubEnvVars()

			if err != nil && tc.err.Error() != err.Error() {
				t.Fatalf("expected %v, got %v", tc.err, err)
			}

			if err == nil {
				if d.localDeploy.DockerHubUserNameEnv != tc.dockerHubUserName {
					t.Fatalf("expected DockerHubUserNameEnv to be %q, got %q", tc.dockerHubUserName, d.localDeploy.DockerHubUserNameEnv)
				}

				if d.localDeploy.DockerHubAccessTokenEnv != tc.dockerHubAccessToken {
					t.Fatalf("expected DockerHubAccessTokenEnv to be %q, got %q", tc.dockerHubAccessToken, d.localDeploy.DockerHubAccessTokenEnv)
				}
			}
		})
	}
}

func TestValidateLocalDeploymentConfig(t *testing.T) {
	tests := []struct {
		desc                 string
		dockerHubUserName    string
		dockerHubAccessToken string
		localDeployment      bool
	}{
		{
			desc:                 "Neither DockerHub flags are present",
			dockerHubUserName:    "",
			dockerHubAccessToken: "",
			localDeployment:      false,
		},
		{
			desc:                 "BothDockerHubAccessToken and DockerHubUserName are specified",
			dockerHubUserName:    "my-docker-hub-username",
			dockerHubAccessToken: "my-docker-hub-access-token",
			localDeployment:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := config.NewInMemoryConfig()
			d := &Deploy{
				config: cfg,
			}
			d.flags.DockerHubUserName = tc.dockerHubUserName
			d.flags.DockerHubAccessToken = tc.dockerHubAccessToken
			d.config.Set(dockerHubUserNameEnv, tc.dockerHubUserName)
			d.config.Set(dockerHubAccessTokenEnv, tc.dockerHubAccessToken)

			err := d.validateLocalDeploymentConfig()

			if err == nil && d.localDeploy.Enabled != tc.localDeployment {
				t.Fatalf("expected localDeployment to be %v, got %v", tc.localDeployment, d.localDeploy.Enabled)
			}
		})
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
	pipelineUUID := "5d0c9667-1626-4ffd-9a94-fab4092eec5a"
	gitSha := "626de930-67ee-4f2b-9af3-12e7165c86b3"

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
		GitSha:   gitSha,
		Pipeline: meroxa.EntityIdentifier{UUID: null.StringFrom(pipelineUUID)},
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

	d.appName = name

	err = d.createApplication(ctx, pipelineUUID, gitSha)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
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

func TestTearDownExistingResourcesWithAppRunning(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	d := &Deploy{
		client: client,
		logger: logger,
	}

	app := utils.GenerateApplication("")
	d.appName = app.Name

	client.
		EXPECT().
		GetApplication(ctx, app.Name).
		Return(&app, nil)

	err := d.tearDownExistingResources(ctx)

	expectedError := fmt.Sprintf("application %q is already %s", app.Name, app.Status.State)
	expectedError = fmt.Sprintf("%s\n\t. Use `meroxa apps remove %s` if you want to redeploy to this application", expectedError, app.Name)
	if err.Error() != expectedError {
		t.Fatalf("expected %s error, got \"%s\"", expectedError, err.Error())
	}
}

func TestTearDownExistingResourcesWithAppDegraded(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	d := &Deploy{
		client: client,
		logger: logger,
	}

	app := utils.GenerateApplication(meroxa.ApplicationStateDegraded)
	d.appName = app.Name

	client.
		EXPECT().
		GetApplication(ctx, app.Name).
		Return(&app, nil)

	res := &http.Response{
		StatusCode: http.StatusNoContent,
	}
	client.
		EXPECT().
		DeleteApplicationEntities(ctx, d.appName).
		Return(res, nil)

	err := d.tearDownExistingResources(ctx)

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}
}

func TestTearDownExistingResourcesWithAppNotFound(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	d := &Deploy{
		client: client,
		logger: logger,
	}

	d.appName = "test"

	client.
		EXPECT().
		GetApplication(ctx, d.appName).
		Return(nil, nil)

	res := &http.Response{
		StatusCode: http.StatusNoContent,
	}
	client.
		EXPECT().
		DeleteApplicationEntities(ctx, d.appName).
		Return(res, nil)

	err := d.tearDownExistingResources(ctx)

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}
}

func TestGetResourceCheckErrorMessage(t *testing.T) {
	testCases := []struct {
		name                 string
		resourceNames        []string
		resourceState        string
		expectedErrorMessage string
	}{
		{
			name:                 "getResourceCheckErrorMessage returns an empty response if all resources are found and available",
			resourceNames:        []string{"nozzle", "engine"},
			resourceState:        "ready",
			expectedErrorMessage: "",
		},
		{
			name:                 "getResourceCheckErrorMessage returns an error response if resources are unavailable",
			resourceNames:        []string{"nozzle", "engine"},
			resourceState:        "",
			expectedErrorMessage: "resource \"nozzle\" is not ready and usable; resource \"engine\" is not ready and usable",
		},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	d := &Deploy{
		client: client,
		logger: logger,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			firstName := "nozzle"
			secondName := "engine"

			firstResource := utils.GenerateResourceWithNameAndStatus(firstName, tc.resourceState)
			secondResource := utils.GenerateResourceWithNameAndStatus(secondName, tc.resourceState)

			client.
				EXPECT().
				GetResourceByNameOrID(ctx, firstResource.Name).
				Return(&firstResource, nil)

			client.
				EXPECT().
				GetResourceByNameOrID(ctx, secondResource.Name).
				Return(&secondResource, nil)

			result := d.getResourceCheckErrorMessage(ctx, tc.resourceNames)
			assert.Equal(t, tc.expectedErrorMessage, result)
		})
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
