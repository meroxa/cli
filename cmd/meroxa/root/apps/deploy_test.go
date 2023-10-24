package apps

import (
	"context"
	"fmt"
	"testing"

	"github.com/meroxa/turbine-core/pkg/ir"

	basicMock "github.com/meroxa/cli/cmd/meroxa/global/mock"
	turbineMock "github.com/meroxa/cli/cmd/meroxa/turbine/mock"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
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

				if err != nil {
					require.Equal(t, newLangUnsupportedError(lang), err)
				} else {
					require.Equal(t, tc.wantErr, err != nil)
				}
			}
		})
	}
}

func TestGetPlatformImage(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()
	appName := "my-app"
	buildPath := ""
	err := fmt.Errorf("nope")

	tests := []struct {
		name           string
		meroxaClient   func(*gomock.Controller) *basicMock.MockBasicClient
		mockTurbineCLI func(*gomock.Controller) turbine.CLI
		err            error
	}{
		{
			name: "Successfully get platform image ",
			meroxaClient: func(ctrl *gomock.Controller) *basicMock.MockBasicClient {
				client := basicMock.NewMockBasicClient(ctrl)
				// client.EXPECT().CollectionRequest(ctx, "POST", "apps", "", nil, nil, &Application{})
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
			name: "Fail to get platform image ",
			meroxaClient: func(ctrl *gomock.Controller) *basicMock.MockBasicClient {
				client := basicMock.NewMockBasicClient(ctrl)
				// client.EXPECT().CollectionRequest(ctx, "POST", "apps", "", nil, nil, &Application{})
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
		meroxaClient   func(*gomock.Controller) *basicMock.MockBasicClient
		mockTurbineCLI func(*gomock.Controller) turbine.CLI
		err            error
	}{
		{
			name: "Don't build app image when for app with no function",
			meroxaClient: func(ctrl *gomock.Controller) *basicMock.MockBasicClient {
				return basicMock.NewMockBasicClient(ctrl)
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
