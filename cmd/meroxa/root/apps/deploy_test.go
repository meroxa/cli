package apps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meroxa/turbine-core/pkg/ir"

	turbineMock "github.com/meroxa/cli/cmd/meroxa/turbine/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployAppFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "path", required: false},
		{name: "spec", required: false, hidden: true},
		{name: "skip-collection-validation", required: false, hidden: false},
		{name: "verbose", required: false, hidden: true},
		{name: "env", required: false, hidden: false},
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

func Test_validateResource(t *testing.T) {
	ctx := context.Background()

	appResources := []turbine.ApplicationResource{
		{Name: "nozzle"},
		{Name: "engine"},
		{Name: "engine"}, // should be dedupped in all cases.
	}
	r1 := utils.GenerateResourceWithNameAndStatus(appResources[0].Name, "ready")
	r2 := utils.GenerateResourceWithNameAndStatus(appResources[1].Name, "ready")

	mockDeploy := func(ctrl *gomock.Controller, r1, r2 meroxa.Resource) *Deploy {
		client := mock.NewMockClient(ctrl)
		client.EXPECT().GetResourceByNameOrID(ctx, r1.Name).Return(&r1, nil)
		client.EXPECT().GetResourceByNameOrID(ctx, r2.Name).Return(&r2, nil)
		return &Deploy{
			client: client,
			logger: log.NewTestLogger(),
		}
	}

	testCases := []struct {
		name        string
		deploy      func(ctrl *gomock.Controller) *Deploy
		envName     string
		resourceEnv string
		state       string
		wantErr     error
	}{
		{
			name:  "resources are valid",
			state: "ready",
			deploy: func(ctrl *gomock.Controller) *Deploy {
				return mockDeploy(ctrl, r1, r2)
			},
		},
		{
			name: "resources are valid in an env",
			deploy: func(ctrl *gomock.Controller) *Deploy {
				d := mockDeploy(
					ctrl,
					utils.ResourceWithEnvironment(r1, "my-env"),
					utils.ResourceWithEnvironment(r2, "my-env"),
				)
				d.flags.Environment = "my-env"

				return d
			},
		},
		{
			name: "invalid when resources are not available",
			deploy: func(ctrl *gomock.Controller) *Deploy {
				return mockDeploy(
					ctrl,
					utils.GenerateResourceWithNameAndStatus(appResources[0].Name, ""),
					utils.GenerateResourceWithNameAndStatus(appResources[1].Name, ""),
				)
			},
			wantErr: errors.New(`resource "nozzle" is not ready and usable; resource "engine" is not ready and usable`),
		},
		{
			name: "invalid when envs do not match",
			deploy: func(ctrl *gomock.Controller) *Deploy {
				d := mockDeploy(
					ctrl,
					utils.ResourceWithEnvironment(r1, "wrong-env"),
					utils.ResourceWithEnvironment(r2, "wrong-env"),
				)
				d.flags.Environment = "test-env"

				return d
			},
			wantErr: errors.New(`resource "nozzle" is not in app env "test-env", but in "wrong-env"; resource "engine" is not in app env "test-env", but in "wrong-env"`), //nolint:lll
		},
		{
			name: "invalid when env is common and resource in not",
			deploy: func(ctrl *gomock.Controller) *Deploy {
				return mockDeploy(
					ctrl,
					utils.ResourceWithEnvironment(r1, "wrong-env"),
					utils.ResourceWithEnvironment(r2, "wrong-env"),
				)
			},
			wantErr: errors.New(`resource "nozzle" is in "wrong-env", but app is in common; resource "engine" is in "wrong-env", but app is in common`), //nolint:lll
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			d := tc.deploy(ctrl)
			err := d.validateResources(ctx, appResources)
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			client := mock.NewMockClient(ctrl)
			logger := log.NewTestLogger()

			d := &Deploy{
				client: client,
				logger: logger,
			}
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
		name      string
		languages []ir.Lang
		wantErr   bool
	}{
		{
			name:      "Successfully validate golang",
			languages: []ir.Lang{"go", "golang"},
		},
		{
			name:      "Successfully validate javascript",
			languages: []ir.Lang{"js", "javascript", "nodejs"},
		},
		{
			name:      "Successfully validate python",
			languages: []ir.Lang{"py", "python", "python3"},
		},
		{
			name:      "Successfully validate ruby",
			languages: []ir.Lang{"rb", "ruby"},
		},
		{
			name:      "Reject unsupported languages",
			languages: []ir.Lang{"cobol", "crystal", "g@rbAg3"},
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, lang := range tc.languages {
				test := Deploy{lang: lang}

				err := test.validateLanguage()
				// require.Equal(t, tc.wantErr, err != nil)

				if err != nil {
					require.Equal(t, newLangUnsupportedError(lang), err)
				} else {
					require.Equal(t, tc.wantErr, err != nil)
				}
			}
		})
	}
}

//nolint:funlen // this is a test function, splitting it would duplicate code
func TestDeployApp(t *testing.T) {
	ctx := context.Background()
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
		name           string
		meroxaClient   func(*gomock.Controller) meroxa.Client
		mockTurbineCLI func(*gomock.Controller, string) turbine.CLI
		version        string
		err            error
	}{
		{
			name: "Successfully deploy app pre IR",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				return client
			},
			mockTurbineCLI: func(ctrl *gomock.Controller, version string) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					GetDeploymentSpec(ctx, imageName).
					Return(specStr, nil)
				return mockTurbineCLI
			},
			version: "",
			err:     nil,
		},
		{
			name: "Successfully deploy app with IR",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
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
			mockTurbineCLI: func(ctrl *gomock.Controller, version string) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					GetDeploymentSpec(ctx, imageName).
					Return(specStr, nil)

				return mockTurbineCLI
			},
			version: specVersion,
			err:     nil,
		},
		{
			name: "Fail to call Turbine deploy",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				return client
			},
			mockTurbineCLI: func(ctrl *gomock.Controller, version string) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					GetDeploymentSpec(ctx, imageName).
					Return(specStr, err)
				return mockTurbineCLI
			},
			version: specVersion,
			err:     err,
		},
		{
			name: "Fail to create deployment",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
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
			mockTurbineCLI: func(ctrl *gomock.Controller, version string) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					GetDeploymentSpec(ctx, imageName).
					Return(specStr, nil)
				return mockTurbineCLI
			},
			version: specVersion,
			err:     err,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			cfg := config.NewInMemoryConfig()
			cfg.Set(global.UserAccountUUID, accountUUID)
			d := &Deploy{
				client:     tc.meroxaClient(ctrl),
				turbineCLI: tc.mockTurbineCLI(ctrl, tc.version),
				logger:     logger,
				appName:    appName,
				config:     cfg,
			}

			_, err := d.createDeployment(ctx, imageName, gitSha, tc.version)
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
	logger := log.NewTestLogger()
	buildUUID := uuid.NewString()
	sourcePutURL := "http://foo.bar"
	sourceGetURL := "http://foo.bar"
	appName := "my-app"
	buildPath := ""
	err := fmt.Errorf("nope")

	tests := []struct {
		name           string
		meroxaClient   func(*gomock.Controller) meroxa.Client
		mockTurbineCLI func(*gomock.Controller) turbine.CLI
		env            string
		err            error
	}{
		{
			name: "Successfully build image with no env",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				input := meroxa.CreateSourceInputV2{}
				client.EXPECT().
					CreateSourceV2(ctx, &input).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: server.URL}, nil)

				client.EXPECT().
					CreateBuild(ctx, &meroxa.CreateBuildInput{SourceBlob: meroxa.SourceBlob{Url: sourceGetURL}}).
					Return(&meroxa.Build{Uuid: buildUUID}, nil)

				client.EXPECT().
					GetBuild(ctx, buildUUID).
					Return(&meroxa.Build{Uuid: buildUUID, Status: meroxa.BuildStatus{State: "complete"}}, nil)
				return client
			},
			mockTurbineCLI: func(ctrl *gomock.Controller) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					CreateDockerfile(ctx, appName).
					Return(buildPath, nil)
				mockTurbineCLI.EXPECT().
					CleanupDockerfile(logger, buildPath).
					Return()
				return mockTurbineCLI
			},
			err: nil,
		},
		{
			name: "Successfully build image with env",
			env:  "my-env",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				input := meroxa.CreateSourceInputV2{Environment: &meroxa.EntityIdentifier{Name: "my-env"}}

				client.EXPECT().
					CreateSourceV2(ctx, &input).
					Return(&meroxa.Source{
						GetUrl:      sourceGetURL,
						PutUrl:      server.URL,
						Environment: &meroxa.EntityIdentifier{Name: "my-env"},
					}, nil)

				client.EXPECT().
					CreateBuild(ctx, &meroxa.CreateBuildInput{
						SourceBlob:  meroxa.SourceBlob{Url: sourceGetURL},
						Environment: &meroxa.EntityIdentifier{Name: "my-env"},
					}).
					Return(&meroxa.Build{Uuid: buildUUID, Environment: &meroxa.EntityIdentifier{Name: "my-env"}}, nil)

				client.EXPECT().
					GetBuild(ctx, buildUUID).
					Return(&meroxa.Build{
						Uuid:   buildUUID,
						Status: meroxa.BuildStatus{State: "complete"},
						Environment: &meroxa.EntityIdentifier{
							Name: "my-env",
						},
					}, nil)
				return client
			},
			mockTurbineCLI: func(ctrl *gomock.Controller) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					CreateDockerfile(ctx, appName).
					Return(buildPath, nil)
				mockTurbineCLI.EXPECT().
					CleanupDockerfile(logger, buildPath).
					Return()
				return mockTurbineCLI
			},
			err: nil,
		},
		{
			name: "Fail to create source",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				input := meroxa.CreateSourceInputV2{}
				client.EXPECT().
					CreateSourceV2(ctx, &input).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: sourcePutURL}, err)
				return client
			},
			mockTurbineCLI: func(ctrl *gomock.Controller) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				return mockTurbineCLI
			},
			err: err,
		},
		{
			name: "Fail to upload source",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				input := meroxa.CreateSourceInputV2{}
				client.EXPECT().
					CreateSourceV2(ctx, &input).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: sourcePutURL}, nil)
				return client
			},
			mockTurbineCLI: func(ctrl *gomock.Controller) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					CreateDockerfile(ctx, appName).
					Return("", err)
				return mockTurbineCLI
			},
			err: err,
		},
		{
			name: "Fail to create build",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				input := meroxa.CreateSourceInputV2{}
				client.EXPECT().
					CreateSourceV2(ctx, &input).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: server.URL}, nil)

				client.EXPECT().
					CreateBuild(ctx, &meroxa.CreateBuildInput{SourceBlob: meroxa.SourceBlob{Url: sourceGetURL}}).
					Return(&meroxa.Build{Uuid: buildUUID}, err)
				return client
			},
			mockTurbineCLI: func(ctrl *gomock.Controller) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					CreateDockerfile(ctx, appName).
					Return(buildPath, nil)
				mockTurbineCLI.EXPECT().
					CleanupDockerfile(logger, buildPath).
					Return()
				return mockTurbineCLI
			},
			err: err,
		},
		{
			name: "Fail to build",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				input := meroxa.CreateSourceInputV2{}
				client.EXPECT().
					CreateSourceV2(ctx, &input).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: server.URL}, nil)

				client.EXPECT().
					CreateBuild(ctx, &meroxa.CreateBuildInput{SourceBlob: meroxa.SourceBlob{Url: sourceGetURL}}).
					Return(&meroxa.Build{Uuid: buildUUID}, nil)

				client.EXPECT().
					GetBuild(ctx, buildUUID).
					Return(&meroxa.Build{Uuid: buildUUID, Status: meroxa.BuildStatus{State: "error"}}, nil)
				return client
			},
			mockTurbineCLI: func(ctrl *gomock.Controller) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					CreateDockerfile(ctx, appName).
					Return(buildPath, nil)
				mockTurbineCLI.EXPECT().
					CleanupDockerfile(logger, buildPath).
					Return()
				return mockTurbineCLI
			},
			err: fmt.Errorf("build with uuid %q errored", buildUUID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			d := &Deploy{
				client:     tc.meroxaClient(ctrl),
				turbineCLI: tc.mockTurbineCLI(ctrl),
				logger:     logger,
				appName:    appName,
			}
			if tc.env != "" {
				d.env = &environment{Name: tc.env}
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

func TestGetAppImage(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()
	appName := "my-app"

	tests := []struct {
		name           string
		meroxaClient   func(*gomock.Controller) meroxa.Client
		mockTurbineCLI func(*gomock.Controller) turbine.CLI
		err            error
	}{
		{
			name: "Don't build app image when for app with no function",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				return mock.NewMockClient(ctrl)
			},
			mockTurbineCLI: func(ctrl *gomock.Controller) turbine.CLI {
				mockTurbineCLI := turbineMock.NewMockCLI(ctrl)
				mockTurbineCLI.EXPECT().
					NeedsToBuild(ctx).
					Return(false, nil)
				return mockTurbineCLI
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			d := &Deploy{
				client:     tc.meroxaClient(ctrl),
				turbineCLI: tc.mockTurbineCLI(ctrl),
				logger:     logger,
				appName:    appName,
			}
			d.flags.Environment = "my-env"

			_, err := d.getAppImage(ctx)
			if err != nil {
				require.NotNil(t, tc.err)
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
		name       string
		branchName string
		resultName string
	}{
		{
			name:       "Prepare app name for main",
			branchName: "main",
			resultName: appName,
		},
		{
			name:       "Prepare app name for master",
			branchName: "master",
			resultName: appName,
		},
		{
			name:       "Prepare app name for feature branch without characters to replace",
			branchName: "my-feature-branch",
			resultName: "my-app-my-feature-branch",
		},
		{
			name:       "Prepare app name for feature branch with characters to replace",
			branchName: "My.Feature\\Branch@01[",
			resultName: "my-app-my-feature-branch-01-",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &Deploy{
				gitBranch:     tc.branchName,
				configAppName: appName,
				logger:        logger,
			}

			result := d.prepareAppName(ctx)
			require.Equal(t, tc.resultName, result)
		})
	}
}

//nolint:funlen
func TestWaitForDeployment(t *testing.T) {
	ctx := context.Background()
	appName := "unit-test"
	outputLogs := []string{"just getting started", "going well", "nailed it"}
	deploymentUuuid := "does-not-matter"

	tests := []struct {
		name         string
		meroxaClient func(*gomock.Controller) meroxa.Client
		wantOutput   func() string
		verboseFlag  bool
		err          error
	}{
		{
			name: "Deployment finishes successfully immediately (no verbose flag)",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeployed,
							Details: strings.Join(outputLogs, "\n"),
						},
					}, nil)
				return client
			},
			wantOutput: func() string { return "" },
			err:        nil,
		},
		{
			name: "Deployment finishes successfully immediately (with verbose flag)",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeployed,
							Details: strings.Join(outputLogs, "\n"),
						},
					}, nil)
				return client
			},
			wantOutput:  func() string { return "\tnailed it\n" },
			verboseFlag: true,
			err:         nil,
		},
		{
			name: "Deployment finishes successfully over time (no verbose flag)",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				first := client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeploying,
							Details: strings.Join(outputLogs[:1], "\n"),
						},
					}, nil)
				second := client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeploying,
							Details: strings.Join(outputLogs[:2], "\n"),
						},
					}, nil)
				third := client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeployed,
							Details: strings.Join(outputLogs, "\n"),
						},
					}, nil).AnyTimes()
				gomock.InOrder(first, second, third)
				return client
			},
			err:        nil,
			wantOutput: func() string { return "" },
		},
		{
			name: "Deployment finishes successfully over time (with verbose flag)",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				first := client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeploying,
							Details: strings.Join(outputLogs[:1], "\n"),
						},
					}, nil)
				second := client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeploying,
							Details: strings.Join(outputLogs[:2], "\n"),
						},
					}, nil)
				third := client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
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
			wantOutput: func() string {
				errorMsg := ""
				for _, l := range outputLogs {
					errorMsg = fmt.Sprintf("%s\t%s\n", errorMsg, l)
				}
				return errorMsg
			},
			verboseFlag: true,
		},
		{
			name: "Deployment immediately failed (no verbose flag)",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeployingError,
							Details: strings.Join(outputLogs, "\n"),
						},
					}, nil)
				return client
			},
			wantOutput: func() string {
				errorMsg := "\n"
				for _, l := range outputLogs {
					errorMsg = fmt.Sprintf("%s\t%s\n", errorMsg, l)
				}
				return errorMsg
			},
			err: fmt.Errorf("\n Check `meroxa apps logs` for further information"),
		},
		{
			name: "Deployment immediately failed (with verbose flag)",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{
						Status: meroxa.DeploymentStatus{
							State:   meroxa.DeploymentStateDeployingError,
							Details: strings.Join(outputLogs, "\n"),
						},
					}, nil)
				return client
			},
			wantOutput: func() string {
				return "\tnailed it\n"
			},
			verboseFlag: true,
			err:         fmt.Errorf("\n Check `meroxa apps logs` for further information"),
		},
		{
			name: "Failed to get latest deployment",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				client.EXPECT().
					GetDeployment(ctx, appName, deploymentUuuid).
					Return(&meroxa.Deployment{}, fmt.Errorf("not today"))
				return client
			},
			wantOutput: func() string { return "" },
			err:        errors.New("couldn't fetch deployment status: not today"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := log.NewTestLogger()
			d := &Deploy{
				client:  tc.meroxaClient(ctrl),
				logger:  logger,
				appName: appName,
			}

			d.flags.Verbose = tc.verboseFlag

			err := d.waitForDeployment(ctx, deploymentUuuid)
			require.Equal(t, tc.err, err, "errors are not equal")

			if err != nil {
				require.Equal(t, tc.wantOutput(), logger.LeveledOutput(), "logs are not equal")
			} else {
				require.Equal(t, tc.wantOutput(), logger.LeveledOutput(), "logs are not equal")
			}
		})
	}
}

func Test_createApplication(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()
	appName := "my-app"

	err := fmt.Errorf(
		`application %q exists in the "running" state.`+"\n"+
			"\tUse 'meroxa apps remove %s' if you want to redeploy to this application",
		appName,
		appName,
	)

	tests := []struct {
		name         string
		meroxaClient func(*gomock.Controller) meroxa.Client
		err          error
	}{
		{
			name: "No need to teardown",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().
					GetApplication(ctx, appName).
					Return(nil, nil)
				client.EXPECT().
					CreateApplicationV2(ctx, &meroxa.CreateApplicationInput{Name: appName}).
					Return(nil, nil)
				return client
			},
			err: nil,
		},
		{
			name: "No need to teardown running app",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().
					GetApplication(ctx, appName).
					Return(&meroxa.Application{Status: meroxa.ApplicationStatus{
						State: meroxa.ApplicationStateRunning,
					}}, nil)
				return client
			},
			err: err,
		},
		{
			name: "Teardown failed app",
			meroxaClient: func(ctrl *gomock.Controller) meroxa.Client {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().
					GetApplication(ctx, appName).
					Return(&meroxa.Application{Status: meroxa.ApplicationStatus{
						State: meroxa.ApplicationStateFailed,
					}}, nil)
				client.EXPECT().
					DeleteApplicationEntities(ctx, appName).
					Return(nil, nil)
				client.EXPECT().
					CreateApplicationV2(ctx, &meroxa.CreateApplicationInput{Name: appName}).
					Return(nil, nil)
				return client
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// cfg := config.NewInMemoryConfig()
			// cfg.Set(global.UserAccountUUID, accountUUID)
			d := &Deploy{
				client:  tc.meroxaClient(ctrl),
				logger:  logger,
				appName: appName,
				//	config:     cfg,
			}

			_, err := d.createApplication(ctx)
			if err != nil {
				require.NotEmpty(t, tc.err)
				require.Equal(t, tc.err, err)
			} else {
				require.Empty(t, tc.err)
			}
		})
	}
}

func TestDeploy_getAppSource(t *testing.T) {
	ctx := context.Background()
	sourceGetURL := "http://foo.bar"

	tests := []struct {
		name         string
		envFlag      string
		meroxaClient func(*gomock.Controller, string) meroxa.Client
	}{
		{
			name:    "when deploying with an environment",
			envFlag: "my-env",
			meroxaClient: func(ctrl *gomock.Controller, env string) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))

				input := meroxa.CreateSourceInputV2{Environment: &meroxa.EntityIdentifier{Name: env}}
				client.EXPECT().
					CreateSourceV2(ctx, &input).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: server.URL}, nil)

				return client
			},
		},
		{
			name: "when deploying without an environment",
			meroxaClient: func(ctrl *gomock.Controller, env string) meroxa.Client {
				client := mock.NewMockClient(ctrl)

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				input := meroxa.CreateSourceInputV2{}

				client.EXPECT().
					CreateSourceV2(ctx, &input).
					Return(&meroxa.Source{GetUrl: sourceGetURL, PutUrl: server.URL}, nil)

				return client
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			d := &Deploy{
				client: tc.meroxaClient(ctrl, tc.envFlag),
			}
			if tc.envFlag != "" {
				d.flags.Environment = tc.envFlag
				d.env = &environment{Name: tc.envFlag}
			}

			s, err := d.getAppSource(ctx)
			require.NoError(t, err)
			assert.NotEmpty(t, s.GetUrl)
			assert.NotEmpty(t, s.PutUrl)
		})
	}
}

func Test_envFromFlag(t *testing.T) {
	tests := []struct {
		desc             string
		flag, uuid, name string
	}{
		{
			desc: "with uuid",
			flag: "543d036e-56af-4ef9-b0a0-f9c55cffac0e",
			uuid: "543d036e-56af-4ef9-b0a0-f9c55cffac0e",
		},
		{
			desc: "with name",
			flag: "env-name",
			name: "env-name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			e := envFromFlag(tc.flag)
			assert.Equal(t, e.Name, tc.name)
			assert.Equal(t, e.UUID, tc.uuid)
		})
	}
}

func Test_validateEnvExists(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		desc    string
		setup   func(ctrl *gomock.Controller) *Deploy
		wantErr error
	}{
		{
			desc: "environment is found",
			setup: func(ctrl *gomock.Controller) *Deploy {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().GetEnvironment(ctx, "my-env").Return(nil, nil)
				d := &Deploy{
					client: client,
					env:    &environment{Name: "my-env"},
				}
				d.flags.Environment = "my-env"
				return d
			},
		},
		{
			desc: "environment is not found",
			setup: func(ctrl *gomock.Controller) *Deploy {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().GetEnvironment(ctx, "your-env").Return(nil,
					fmt.Errorf("could not find environment"),
				)
				d := &Deploy{
					client: client,
					env:    &environment{Name: "your-env"},
				}
				d.flags.Environment = "your-env"
				return d
			},
			wantErr: fmt.Errorf(`environment "your-env" does not exist`),
		},
		{
			desc: "failed to retrieve environment",
			setup: func(ctrl *gomock.Controller) *Deploy {
				client := mock.NewMockClient(ctrl)
				client.EXPECT().GetEnvironment(ctx, "your-env").Return(nil,
					fmt.Errorf("boom"),
				)
				d := &Deploy{
					client: client,
					env:    &environment{Name: "your-env"},
				}
				d.flags.Environment = "your-env"
				return d
			},
			wantErr: fmt.Errorf(`unable to retrieve environment "your-env": boom`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			d := tc.setup(gomock.NewController(t))
			err := d.validateEnvExists(ctx)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
