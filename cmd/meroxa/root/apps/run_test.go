package apps

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	mockturbinecli "github.com/meroxa/cli/cmd/meroxa/turbine/mock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
)

func TestRunAppFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "path", required: false},
	}

	c := builder.BuildCobraCommand(&Run{})

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

func TestRunExecute(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCli := func() turbine.CLI {
		mock := mockturbinecli.NewMockCLI(mockCtrl)
		mock.EXPECT().Run(ctx)
		return mock
	}

	tests := []struct {
		desc     string
		cli      turbine.CLI
		config   turbine.AppConfig
		features []string
		err      error
	}{
		{
			desc: "Execute Javascript run successfully",
			cli:  mockCli(),
			config: turbine.AppConfig{
				Name:     "js-test",
				Language: turbine.JavaScript,
				Vendor:   "false",
			},
		},
		{
			desc: "Execute Golang run successfully",
			cli:  mockCli(),
			config: turbine.AppConfig{
				Name:     "go-test",
				Language: turbine.GoLang,
			},
		},
		{
			desc: "Execute Ruby Run successfully",
			cli:  mockCli(),
			config: turbine.AppConfig{
				Name:     "ruby-test",
				Language: turbine.Ruby,
			},
			features: []string{"ruby_implementation"},
		},
		{
			desc: "Execute Ruby Run (Missing feature)",
			cli:  mockturbinecli.NewMockCLI(mockCtrl),
			config: turbine.AppConfig{
				Name:     "ruby-test",
				Language: turbine.Ruby,
			},
			err: fmt.Errorf(`no access to the Meroxa Turbine Ruby feature.
Sign up for the Beta here: https://share.hsforms.com/1Uq6UYoL8Q6eV5QzSiyIQkAc2sme`),
		},

		{
			desc: "Execute Python Run with an error",
			cli: func() turbine.CLI {
				mock := mockturbinecli.NewMockCLI(mockCtrl)
				mock.EXPECT().Run(ctx).Return(fmt.Errorf("not good"))
				return mock
			}(),
			config: turbine.AppConfig{
				Name:     "py-test",
				Language: turbine.Python,
			},
			err: fmt.Errorf("not good"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			u := &Run{
				logger:     log.NewTestLogger(),
				config:     &tt.config,
				turbineCLI: tt.cli,
			}

			if global.Config == nil {
				build := builder.BuildCobraCommand(u)
				_ = global.PersistentPreRunE(build)
			}

			if len(tt.features) != 0 {
				oldflags := setFlags(tt.features, false)
				defer setFlags(oldflags, true)
			}

			err := u.Execute(ctx)
			processError(t, err, tt.err)
			if err == nil && tt.err != nil {
				t.Fatalf("did not find expected error: %s", tt.err.Error())
			}
		})
	}
}

// setFlags adds newflags to the global flag collection and returns old flags,
// when replace is true, flags will be overwritten with newflags.
func setFlags(newflags []string, replace bool) []string {
	var flags string
	oldflags := global.Config.Get(global.UserFeatureFlagsEnv)
	if oldflags != nil {
		flags = oldflags.(string)
	} else {
		oldflags = ""
	}

	if replace {
		flags = strings.Join(newflags, " ")
	} else {
		flags += " " + strings.Join(newflags, " ")
	}
	global.Config.Set(global.UserFeatureFlagsEnv, flags)

	return strings.Split(oldflags.(string), " ")
}
