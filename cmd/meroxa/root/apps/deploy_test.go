package apps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	turbine_mock "github.com/meroxa/cli/cmd/meroxa/turbine/mock"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
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
		{name: "spec", required: false, hidden: true},
		{name: "skip-collection-validation", required: false, hidden: false},
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
		resources            []turbine.ApplicationResource
		resourceState        string
		expectedErrorMessage string
	}{
		{
			name: "getResourceCheckErrorMessage returns an empty response if all resources are found and available",
			resources: []turbine.ApplicationResource{
				{
					Name: "nozzle",
				},
				{
					Name: "engine",
				},
			},
			resourceState:        "ready",
			expectedErrorMessage: "",
		},
		{
			name: "getResourceCheckErrorMessage returns an error response if resources are unavailable",
			resources: []turbine.ApplicationResource{
				{
					Name: "nozzle",
				},
				{
					Name: "engine",
				},
			},
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

			result := d.getResourceCheckErrorMessage(ctx, tc.resources)
			assert.Equal(t, tc.expectedErrorMessage, result)
		})
	}
}

//nolint:funlen // this is a test function, splitting it would duplicate code
func TestValidateCollections(t *testing.T) {
	testCases := []struct {
		name      string
		resources []turbine.ApplicationResource
		err       string
	}{
		{
			name: "Different source and destination resources reference different collections",
			resources: []turbine.ApplicationResource{
				{
					Name:       "source",
					Source:     true,
					Collection: "sequences",
				},
				{
					Name:        "destination",
					Destination: true,
					Collection:  "test-destination",
				},
			},
		},
		{
			name: "Different source and destination resources reference same collection",
			resources: []turbine.ApplicationResource{
				{
					Name:       "source",
					Source:     true,
					Collection: "sequences",
				},
				{
					Name:        "destination",
					Destination: true,
					Collection:  "sequences",
				},
			},
		},
		{
			name: "Multiple destination resources",
			resources: []turbine.ApplicationResource{
				{
					Name:       "source",
					Source:     true,
					Collection: "sequences",
				},
				{
					Name:        "destination",
					Destination: true,
					Collection:  "sequences",
				},
				{
					Name:        "alt-destination",
					Destination: true,
					Collection:  "sequences",
				},
			},
		},
		{
			name: "Same source and destination resources reference same collection",
			resources: []turbine.ApplicationResource{
				{
					Name:       "pg",
					Source:     true,
					Collection: "sequences",
				},
				{
					Name:        "pg",
					Destination: true,
					Collection:  "sequences",
				},
			},
			err: "⚠️\n\tApplication resource \"pg\" with collection \"sequences\" cannot be used as a destination. It is also the source." +
				"\nPlease modify your Turbine data application code. Then run `meroxa app deploy` again. " +
				"To skip collection validation, run `meroxa app deploy --skip-collection-validation`.",
		},
		{
			name: "One resource is both source and destination",
			resources: []turbine.ApplicationResource{
				{
					Name:        "source",
					Source:      true,
					Destination: true,
					Collection:  "sequences",
				},
			},
			err: "⚠️\n\tApplication resource cannot be used as both a source and destination." +
				"\nPlease modify your Turbine data application code. Then run `meroxa app deploy` again. " +
				"To skip collection validation, run `meroxa app deploy --skip-collection-validation`.",
		},
		{
			name: "Destination resource used in another app",
			resources: []turbine.ApplicationResource{
				{
					Name:       "source",
					Source:     true,
					Collection: "sequences",
				},
				{
					Name:        "pg",
					Destination: true,
					Collection:  "anonymous",
				},
			},
			err: "⚠️\n\tApplication resource \"pg\" with collection \"anonymous\" cannot be used as a destination. " +
				"It is also being used as a destination by another application \"application-name\"." +
				"\nPlease modify your Turbine data application code. Then run `meroxa app deploy` again. " +
				"To skip collection validation, run `meroxa app deploy --skip-collection-validation`.",
		},
		{
			name: "Two same destination resources",
			resources: []turbine.ApplicationResource{
				{
					Name:       "source",
					Source:     true,
					Collection: "sequences",
				},
				{
					Name:        "pg",
					Destination: true,
					Collection:  "test-destination",
				},
				{
					Name:        "pg",
					Destination: true,
					Collection:  "test-destination",
				},
			},
			err: "⚠️\n\tApplication resource \"pg\" with collection \"test-destination\" cannot be used as a destination more than once." +
				"\nPlease modify your Turbine data application code. Then run `meroxa app deploy` again. " +
				"To skip collection validation, run `meroxa app deploy --skip-collection-validation`.",
		},
		{
			name: "Ignore resources without collection info",
			resources: []turbine.ApplicationResource{
				{
					Name: "source",
				},
				{
					Name: "pg",
				},
				{
					Name: "pg",
				},
			},
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
			apps := []*meroxa.Application{
				{
					Name: "application-name",
					Resources: []meroxa.ApplicationResource{
						{
							EntityIdentifier: meroxa.EntityIdentifier{
								Name: "pg",
							},
							Collection: meroxa.ResourceCollection{
								Name:        "anonymous",
								Destination: "true",
							},
						},
						{
							EntityIdentifier: meroxa.EntityIdentifier{
								Name: "source",
							},
							Collection: meroxa.ResourceCollection{
								Name:   "sequences",
								Source: "true",
							},
						},
					},
				},
			}
			client.
				EXPECT().
				ListApplications(ctx).
				Return(apps, nil)

			err := d.validateCollections(ctx, tc.resources)
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err.Error(), tc.err)
			}
		})
	}
}

func TestValidateLanguage(t *testing.T) {
	tests := []struct {
		description string
		languages   []string
		errFormat   string
	}{
		{
			description: "Successfully validate golang",
			languages:   []string{"go", "golang"},
		},
		{
			description: "Successfully validate javascript",
			languages:   []string{"js", "javascript", "nodejs"},
		},
		{
			description: "Successfully validate python",
			languages:   []string{"py", "python", "python3"},
		},
		{
			description: "Reject unsupported languages",
			languages:   []string{"cobol", "crystal", "g@rbAg3"},
			errFormat:   "language %q not supported. Currently, we support \"javascript\", \"golang\", and \"python\"",
		},
	}

	test := &Deploy{}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			for _, lang := range tc.languages {
				test.lang = lang
				err := test.validateLanguage()

				if err != nil {
					require.NotEmptyf(t, tc.errFormat, fmt.Sprintf("test failed for %q", lang))
					require.Equal(t, fmt.Errorf(tc.errFormat, lang), err)
				} else {
					require.Emptyf(t, tc.errFormat, "got an unexpected error for input "+lang)
				}
			}
		})
	}
}

//nolint:funlen // this is a test function, splitting it would duplicate code
func TestDeployApp(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	logger := log.NewTestLogger()
	appName := "my-app"
	imageName := "doc.ker:latest"
	gitSha := "aa:bb:cc:dd"
	specVersion := "latest"
	accountUUID := "aa-bb-cc-dd"
	specStr := `{"metadata": "metadata"}`
	spec := map[string]interface{}{
		"metadata": "metadata",
	}
	err := fmt.Errorf("nope")

	tests := []struct {
		description    string
		meroxaClient   func() meroxa.Client
		mockTurbineCLI func(string) turbine.CLI
		version        string
		err            error
	}{
		{
			description: "Successfully deploy app pre IR",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)
				return client
			},
			mockTurbineCLI: func(version string) turbine.CLI {
				mockTurbineCLI := turbine_mock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					Deploy(ctx, imageName, appName, gitSha, version, accountUUID).
					Return(specStr, nil)
				return mockTurbineCLI
			},
			version: "",
			err:     nil,
		},
		{
			description: "Successfully deploy app with IR",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)
				input := &meroxa.CreateDeploymentInput{
					Application: meroxa.EntityIdentifier{Name: appName},
					GitSha:      gitSha,
					SpecVersion: specVersion,
					Spec:        spec,
				}
				client.EXPECT().
					CreateDeployment(ctx, input).
					Return(&meroxa.Deployment{}, nil)
				return client
			},
			mockTurbineCLI: func(version string) turbine.CLI {
				mockTurbineCLI := turbine_mock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					Deploy(ctx, imageName, appName, gitSha, version, accountUUID).
					Return(specStr, nil)

				return mockTurbineCLI
			},
			version: specVersion,
			err:     nil,
		},
		{
			description: "Fail to call Turbine deploy",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)
				return client
			},
			mockTurbineCLI: func(version string) turbine.CLI {
				mockTurbineCLI := turbine_mock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					Deploy(ctx, imageName, appName, gitSha, version, accountUUID).
					Return(specStr, err)
				return mockTurbineCLI
			},
			version: specVersion,
			err:     err,
		},
		{
			description: "Fail to create deployment",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					CreateDeployment(ctx, &meroxa.CreateDeploymentInput{
						Application: meroxa.EntityIdentifier{Name: appName},
						GitSha:      gitSha,
						SpecVersion: specVersion,
						Spec:        spec,
					}).
					Return(&meroxa.Deployment{}, err)
				return client
			},
			mockTurbineCLI: func(version string) turbine.CLI {
				mockTurbineCLI := turbine_mock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					Deploy(ctx, imageName, appName, gitSha, version, accountUUID).
					Return(specStr, nil)
				return mockTurbineCLI
			},
			version: specVersion,
			err:     err,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			cfg := config.NewInMemoryConfig()
			cfg.Set(global.UserAccountUUID, accountUUID)
			d := &Deploy{
				client:     tc.meroxaClient(),
				turbineCLI: tc.mockTurbineCLI(tc.version),
				logger:     logger,
				appName:    appName,
				config:     cfg,
			}

			_, err := d.deployApp(ctx, imageName, gitSha, tc.version)
			if err != nil {
				require.NotEmpty(t, tc.err)
				require.Equal(t, tc.err, err)
			} else {
				require.Empty(t, tc.err)
			}
		})
	}
}

//nolint:funlen // this is a test function, splitting it would duplicate code
func TestGetPlatformImage(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	logger := log.NewTestLogger()
	buildUUID := uuid.NewString()
	sourcePutURL := "http://foo.bar"
	sourceGetURL := "http://foo.bar"
	appName := "my-app"
	tempPath := "/tmp"
	err := fmt.Errorf("nope")

	tests := []struct {
		description    string
		meroxaClient   func() meroxa.Client
		mockTurbineCLI func() turbine.CLI
		err            error
	}{
		{
			description: "Successfully build image",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					CreateSource(ctx).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: sourcePutURL}, nil)

				client.EXPECT().
					CreateBuild(ctx, &meroxa.CreateBuildInput{SourceBlob: meroxa.SourceBlob{Url: sourceGetURL}}).
					Return(&meroxa.Build{Uuid: buildUUID}, nil)

				client.EXPECT().
					GetBuild(ctx, buildUUID).
					Return(&meroxa.Build{Uuid: buildUUID, Status: meroxa.BuildStatus{State: "complete"}}, nil)
				return client
			},
			mockTurbineCLI: func() turbine.CLI {
				mockTurbineCLI := turbine_mock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					UploadSource(ctx, appName, tempPath, sourcePutURL).
					Return(nil)
				return mockTurbineCLI
			},
			err: nil,
		},
		{
			description: "Fail to create source",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					CreateSource(ctx).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: sourcePutURL}, err)
				return client
			},
			mockTurbineCLI: func() turbine.CLI {
				mockTurbineCLI := turbine_mock.NewMockCLI(ctrl)
				return mockTurbineCLI
			},
			err: err,
		},
		{
			description: "Fail to upload source",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					CreateSource(ctx).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: sourcePutURL}, nil)
				return client
			},
			mockTurbineCLI: func() turbine.CLI {
				mockTurbineCLI := turbine_mock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					UploadSource(ctx, appName, tempPath, sourcePutURL).
					Return(err)
				return mockTurbineCLI
			},
			err: err,
		},
		{
			description: "Fail to create build",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					CreateSource(ctx).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: sourcePutURL}, nil)

				client.EXPECT().
					CreateBuild(ctx, &meroxa.CreateBuildInput{SourceBlob: meroxa.SourceBlob{Url: sourceGetURL}}).
					Return(&meroxa.Build{Uuid: buildUUID}, err)
				return client
			},
			mockTurbineCLI: func() turbine.CLI {
				mockTurbineCLI := turbine_mock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					UploadSource(ctx, appName, tempPath, sourcePutURL).
					Return(nil)
				return mockTurbineCLI
			},
			err: err,
		},
		{
			description: "Fail to build",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					CreateSource(ctx).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: sourcePutURL}, nil)

				client.EXPECT().
					CreateBuild(ctx, &meroxa.CreateBuildInput{SourceBlob: meroxa.SourceBlob{Url: sourceGetURL}}).
					Return(&meroxa.Build{Uuid: buildUUID}, nil)

				client.EXPECT().
					GetBuild(ctx, buildUUID).
					Return(&meroxa.Build{Uuid: buildUUID, Status: meroxa.BuildStatus{State: "error"}}, nil)
				return client
			},
			mockTurbineCLI: func() turbine.CLI {
				mockTurbineCLI := turbine_mock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					UploadSource(ctx, appName, tempPath, sourcePutURL).
					Return(nil)
				return mockTurbineCLI
			},
			err: fmt.Errorf("build with uuid %q errored", buildUUID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			d := &Deploy{
				client:     tc.meroxaClient(),
				turbineCLI: tc.mockTurbineCLI(),
				logger:     logger,
				appName:    appName,
				tempPath:   tempPath,
			}

			_, err := d.getPlatformImage(ctx)
			if err != nil {
				require.NotEmpty(t, tc.err)
				require.Equal(t, tc.err, err)
			} else {
				require.Empty(t, tc.err)
			}
		})
	}
}

func TestPrepareAppName(t *testing.T) {
	ctx := context.Background()
	appName := "my-app"
	logger := log.NewTestLogger()

	tests := []struct {
		description string
		branchName  string
		resultName  string
	}{
		{
			description: "Prepare app name for main",
			branchName:  "main",
			resultName:  appName,
		},
		{
			description: "Prepare app name for master",
			branchName:  "master",
			resultName:  appName,
		},
		{
			description: "Prepare app name for feature branch without characters to replace",
			branchName:  "my-feature-branch",
			resultName:  "my-app-my-feature-branch",
		},
		{
			description: "Prepare app name for feature branch with characters to replace",
			branchName:  "My.Feature\\Branch@01[",
			resultName:  "my-app-my-feature-branch-01-",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			d := &Deploy{
				gitBranch:     tc.branchName,
				configAppName: appName,
				logger:        logger,
			}

			result := d.prepareAppName(ctx)
			fmt.Printf("result: %s\n", result)
			require.Equal(t, tc.resultName, result)
		})
	}
}

//nolint:funlen
func TestWaitForDeployment(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	appName := "unit-test"
	outputLogs := []string{"just getting started", "going well", "nailed it"}
	uuid := "does=not-matter"

	tests := []struct {
		description  string
		meroxaClient func() meroxa.Client
		noLogs       bool
		err          error
	}{
		{
			description: "Deployment immediately finished",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					GetDeployment(ctx, appName, uuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeployed,
							Details: strings.Join(outputLogs, "\n"),
						},
					}, nil)
				return client
			},
			err: nil,
		},
		{
			description: "Deployment finishes over time",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				first := client.EXPECT().
					GetDeployment(ctx, appName, uuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeploying,
							Details: strings.Join(outputLogs[:1], "\n"),
						},
					}, nil)
				second := client.EXPECT().
					GetDeployment(ctx, appName, uuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeploying,
							Details: strings.Join(outputLogs[:2], "\n"),
						},
					}, nil)
				third := client.EXPECT().
					GetDeployment(ctx, appName, uuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeployed,
							Details: strings.Join(outputLogs, "\n"),
						},
					}, nil).AnyTimes()
				gomock.InOrder(first, second, third)
				return client
			},
			err: nil,
		},
		{
			description: "Deployment immediately failed",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					GetDeployment(ctx, appName, uuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateRollingBackError,
							Details: strings.Join(outputLogs, "\n"),
						},
					}, nil)
				return client
			},
			noLogs: false,
			err:    fmt.Errorf("failed to deploy Application %q", appName),
		},
		{
			description: "Failed to get latest deployment",
			meroxaClient: func() meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					GetDeployment(ctx, appName, uuid).
					Return(&meroxa.Deployment{}, fmt.Errorf("not today"))
				return client
			},
			noLogs: true,
			err:    fmt.Errorf("not today"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			logger := log.NewTestLogger()
			d := &Deploy{
				client:  tc.meroxaClient(),
				logger:  logger,
				appName: appName,
			}

			err := d.waitForDeployment(ctx, uuid)
			require.Equal(t, tc.err, err, "errors are not equal")

			if err != nil {
				var failureLogs string
				if tc.noLogs {
					failureLogs = ""
				} else {
					failureLogs = outputLogs[0] + "\n" + outputLogs[1] + "\n\tx " + outputLogs[2] + "\n"
				}
				require.Equal(t, failureLogs, logger.SpinnerOutput(), "logs are not equal")
			} else {
				successLogs := strings.Join(outputLogs, "\n") + "\n"
				require.Equal(t, successLogs, logger.SpinnerOutput(), "logs are not equal")
			}
		})
	}
}
