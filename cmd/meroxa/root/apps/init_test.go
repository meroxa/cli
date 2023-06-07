package apps

import (
	"context"
	"errors"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/meroxa/turbine-core/pkg/ir"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	mockturbinecli "github.com/meroxa/cli/cmd/meroxa/turbine/mock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
)

func TestInitAppArgs(t *testing.T) {
	tests := []struct {
		args    []string
		err     error
		appName string
	}{
		{args: nil, err: errors.New("requires an application name"), appName: ""},
		{args: []string{"my-app-name"}, err: nil, appName: "my-app-name"},
	}

	for _, tt := range tests {
		cc := &Init{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.appName != cc.args.appName {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.appName, cc.args.appName)
		}
	}
}

func TestInitAppFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "lang", shorthand: "l", required: true},
		{name: "path", required: false},
	}

	c := builder.BuildCobraCommand(&Init{})

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

//nolint:funlen
func TestGoInitExecute(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	if gopath == "" {
		t.Fatal("GOPATH is not defined")
	}

	tests := []struct {
		desc                 string
		path                 string
		skipModInit          bool
		effectiveSkipModInit bool
		vendor               bool
		err                  error
	}{
		{
			desc: "Init without go mod init",
			path: func() string {
				return filepath.Join(gopath, "src/github.com/meroxa/tests", uuid.New().String()) //nolint:gocritic
			}(),
			skipModInit: true,
			vendor:      false,
			err:         nil,
		},
		{
			desc: "Init with go mod init and without vendoring",
			path: func() string {
				return filepath.Join(gopath, "src/github.com/meroxa/tests", uuid.New().String()) //nolint:gocritic
			}(),
			skipModInit: false,
			vendor:      false,
			err:         nil,
		},
		{
			desc: "Init with go mod init and with vendoring",
			path: func() string {
				return filepath.Join(gopath, "src/github.com/meroxa/tests", uuid.New().String()) //nolint:gocritic
			}(),
			skipModInit: false,
			vendor:      true,
			err:         nil,
		},
		{
			desc: "Init without go mod init and with vendor flag",
			path: func() string {
				return filepath.Join(gopath, "src/github.com/meroxa/tests", uuid.New().String()) //nolint:gocritic
			}(),
			skipModInit: true,
			vendor:      true,
			err:         nil,
		},
		{
			desc: "Init without go mod init because the path is not under GOPATH and with vendor flag",
			path: func() string {
				return filepath.Join(os.TempDir(), uuid.New().String()) //nolint:gocritic
			}(),
			skipModInit:          false,
			effectiveSkipModInit: true,
			vendor:               true,
			err:                  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cc := &Init{}
			cc.Logger(log.NewTestLogger())
			cc.flags.Path = tt.path
			cc.flags.Lang = string(ir.GoLang)
			cc.flags.ModVendor = tt.vendor
			cc.flags.SkipModInit = tt.skipModInit
			cc.args.appName = "app-name"

			err := cc.Execute(context.Background())
			if err != nil {
				if tt.err == nil {
					t.Fatalf("unexpected error \"%s\"", err)
				} else if tt.err.Error() != err.Error() {
					t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
				}
			}

			if err == nil && tt.err == nil {
				if !tt.skipModInit && !tt.effectiveSkipModInit {
					p, _ := filepath.Abs(filepath.Join(tt.path, cc.args.appName, "go.mod"))
					if _, err := os.Stat(p); os.IsNotExist(err) {
						t.Fatalf("expected file \"%s\" to be created", p)
					}

					if tt.vendor {
						p, _ = filepath.Abs(filepath.Join(tt.path, cc.args.appName, "go.mod"))
						if _, err := os.Stat(p); os.IsNotExist(err) {
							t.Fatalf("expected directory \"%s\" to be created", p)
						}
					}
				}
			}
			os.RemoveAll(tt.path)
		})
	}
}

func TestJsPyAndRbInitExecute(t *testing.T) {
	tests := []struct {
		desc     string
		path     string
		language ir.Lang
		name     string
		err      error
	}{
		{
			desc:     "Successful Javascript init",
			path:     "/does/not/matter",
			language: ir.JavaScript,
			name:     "js-init",
			err:      nil,
		},
		{
			desc:     "Unsuccessful Javascript init",
			path:     "/does/not/matter",
			language: ir.JavaScript,
			name:     "js-init",
			err:      fmt.Errorf("not good"),
		},
		{
			desc:     "Successful Python init",
			path:     "/does/not/matter",
			language: ir.Python,
			name:     "py-init",
			err:      nil,
		},
		{
			desc:     "Unsuccessful Python init",
			path:     "/does/not/matter",
			language: ir.Python,
			name:     "py-init",
			err:      fmt.Errorf("not good"),
		},
		{
			desc:     "Successful Ruby init",
			path:     "/does/not/matter",
			language: ir.Ruby,
			name:     "rb-init",
			err:      nil,
		},
		{
			desc:     "Unsuccessful Ruby init",
			path:     "/does/not/matter",
			language: ir.Ruby,
			name:     "rb-init",
			err:      fmt.Errorf("not good"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctx := context.Background()
			mockCtrl := gomock.NewController(t)

			i := &Init{}
			i.Logger(log.NewTestLogger())
			i.flags.Path = tt.path
			i.args.appName = tt.name
			i.flags.Lang = string(tt.language)

			mock := mockturbinecli.NewMockCLI(mockCtrl)
			if tt.err == nil {
				mock.EXPECT().Init(ctx, tt.name)
				mock.EXPECT().GitInit(ctx, filepath.Join(tt.path, tt.name))
			} else {
				mock.EXPECT().Init(ctx, tt.name).Return(tt.err)
			}
			i.turbineCLI = mock

			err := i.Execute(ctx)
			processError(t, err, tt.err)
			if err == nil && tt.err != nil {
				t.Fatalf("did not find expected error: %s", tt.err.Error())
			}
		})
	}
}

func TestAppNameValidation(t *testing.T) {
	tests := []struct {
		desc       string
		inputName  string
		outputName string
		err        error
	}{
		{
			desc:       "Valid app name",
			inputName:  "perfect-name",
			outputName: "perfect-name",
			err:        nil,
		},
		{
			desc:       "Valid capitalized app name",
			inputName:  "Perfect-name",
			outputName: "perfect-name",
			err:        nil,
		},
		{
			desc:       "Valid app name with underscore",
			inputName:  "perfect_name",
			outputName: "perfect_name",
			err:        nil,
		},
		{
			desc:       "Invalid app name - leading number",
			inputName:  "3otherwisegoodname",
			outputName: "",
			err: fmt.Errorf("invalid application name: %s; should start with a letter, be alphanumeric,"+
				" and only have dashes as separators", "3otherwisegoodname"),
		},
		{
			desc:       "Invalid app name - invalid characters",
			inputName:  "!ch@os",
			outputName: "",
			err: fmt.Errorf("invalid application name: %s; should start with a letter, be alphanumeric,"+
				" and only have dashes as separators", "!ch@os"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			output, err := validateAppName(tt.inputName)
			if err != nil {
				if tt.err == nil {
					t.Fatalf("unexpected error \"%s\"", err)
				} else if tt.err.Error() != err.Error() {
					t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
				}
			}

			if err == nil && tt.err == nil {
				if output != tt.outputName {
					t.Fatalf("expected \"%s\" got \"%s\"", tt.outputName, output)
				}
			}
		})
	}
}
