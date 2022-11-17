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
		features map[string]bool
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
			features: map[string]bool{
				"ruby_implementation": true,
			},
		},
		{
			desc: "Execute Ruby Run (Missing feature)",
			cli:  mockturbinecli.NewMockCLI(mockCtrl),
			config: turbine.AppConfig{
				Name:     "ruby-test",
				Language: turbine.Ruby,
			},
			features: map[string]bool{
				"ruby_implementation": false,
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
				setFeatures(tt.features)
				defer resetFeatures(tt.features)
			}

			err := u.Execute(ctx)
			processError(t, err, tt.err)
			if err == nil && tt.err != nil {
				t.Fatalf("did not find expected error: %s", tt.err.Error())
			}
		})
	}
}

// setFeatures sets features from a map which designates enabled/disabled features.
func setFeatures(features map[string]bool) {
	currentFlags := getFeatures()

	for k, v := range features {
		switch v {
		case true:
			currentFlags[k] = v
		case false:
			delete(currentFlags, k)
		}
	}

	var flags []string
	for k := range currentFlags {
		flags = append(flags, k)
	}

	global.Config.Set(global.UserFeatureFlagsEnv, strings.Join(flags, " "))
}

// resetFeatures inverts the state of features defined in the map.
func resetFeatures(features map[string]bool) {
	reset := make(map[string]bool)
	for k, v := range features {
		reset[k] = !v
	}

	setFeatures(reset)
}

func getFeatures() map[string]bool {
	flags := make(map[string]bool)

	if s := global.Config.Get(global.UserFeatureFlagsEnv); s != nil {
		for _, t := range strings.Split(s.(string), " ") {
			flags[t] = true
		}
	}

	return flags
}
