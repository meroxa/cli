package apps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/root/nop"
	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine"
	mockturbinecli "github.com/meroxa/cli/cmd/meroxa/turbine/mock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
)

func TestUpgradeAppFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "path", required: false},
	}

	c := builder.BuildCobraCommand(&Upgrade{})

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

func TestUpgradeExecute(t *testing.T) {
	tests := []struct {
		desc   string
		path   string
		lang   string
		vendor bool
		err    error
	}{
		{
			desc:   "Golang upgrade without vendor",
			path:   "/tmp",
			lang:   GoLang,
			vendor: false,
			err:    nil,
		},
		{
			desc:   "Golang upgrade with path and vendor",
			path:   "/tmp",
			lang:   GoLang,
			vendor: true,
			err:    nil,
		},
		{
			desc:   "Golang upgrade with path and vendor and error",
			path:   "/tmp",
			lang:   GoLang,
			vendor: true,
			err:    fmt.Errorf("not good"),
		},
		{
			desc: "Javascript with path",
			path: "/tmp",
			lang: JavaScript,
			err:  nil,
		},
		{
			desc: "Python with path",
			path: "/tmp",
			lang: Python,
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_ = os.Setenv("UNIT_TEST", "true")
			defer func() {
				_ = os.Unsetenv("UNIT_TEST")
			}()
			mockCtrl := gomock.NewController(t)

			u := &Upgrade{}
			u.Logger(log.NewTestLogger())
			path, err := turbinecli.GetPath(tt.path)
			processError(t, err, tt.err)
			u.flags.Path = path
			u.run = &nop.Nop{}

			mock := mockturbinecli.NewMockCLI(mockCtrl)
			if tt.err == nil {
				mock.EXPECT().Upgrade(path, tt.vendor)
			} else {
				mock.EXPECT().Upgrade(path, tt.vendor).Return(tt.err)
			}
			u.turbineCLI = mock

			switch tt.lang {
			case GoLang:
				config := turbinecli.AppConfig{Language: GoLang}
				config.Vendor = "false"
				if tt.vendor {
					config.Vendor = "true"
				}
				err = turbinecli.WriteConfigFile(path, config)
				processError(t, err, tt.err)
				defer func() {
					_ = os.Remove(filepath.Join(path, "app.json"))
				}()
			case JavaScript:
				_ = turbinecli.WriteConfigFile(path, turbinecli.AppConfig{Language: JavaScript})
				defer func() {
					_ = os.Remove(filepath.Join(path, "app.json"))
				}()
			case Python:
				_ = turbinecli.WriteConfigFile(path, turbinecli.AppConfig{Language: Python})
				defer func() {
					_ = os.Remove(filepath.Join(path, "app.json"))
				}()
			default:
				t.Fatalf("unprocessable language: %s", tt.lang)
			}

			err = u.Execute(context.Background())
			processError(t, err, tt.err)
			if err == nil && tt.err != nil {
				t.Fatalf("did not find expected error: %s", tt.err.Error())
			}
		})
	}
}

func processError(t *testing.T, given error, wanted error) {
	if given != nil {
		if wanted == nil {
			t.Fatalf("unexpected error \"%s\"", given)
		} else if wanted.Error() != given.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", wanted, given)
		}
	}
}
