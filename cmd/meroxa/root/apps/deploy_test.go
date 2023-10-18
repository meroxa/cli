package apps

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meroxa/turbine-core/pkg/ir"

	turbineMock "github.com/meroxa/cli/cmd/meroxa/turbine/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
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

			err := d.getPlatformImage(ctx)
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

			err := d.getAppImage(ctx)
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
