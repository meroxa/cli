package apps

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/root/nop"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
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
		config turbine.AppConfig
		err    error
	}{
		{
			desc: "Successful Javascript upgrade without vendor",
			config: turbine.AppConfig{
				Name:     "",
				Language: JavaScript,
				Vendor:   "false",
			},
			err: nil,
		},
		{
			desc: "Successful Golang upgrade with vendor",
			config: turbine.AppConfig{
				Name:     "",
				Language: GoLang,
				Vendor:   "true",
			},
			err: nil,
		},
		{
			desc: "Unsuccessful Python upgrade",
			config: turbine.AppConfig{
				Name:     "",
				Language: Python,
				Vendor:   "false",
			},
			err: fmt.Errorf("not good"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)

			u := &Upgrade{}
			u.Logger(log.NewTestLogger())
			u.config = &tt.config
			u.run = &nop.Nop{}
			vendor, _ := strconv.ParseBool(tt.config.Vendor)

			mock := mockturbinecli.NewMockCLI(mockCtrl)
			if tt.err == nil {
				mock.EXPECT().Upgrade(vendor)
			} else {
				mock.EXPECT().Upgrade(vendor).Return(tt.err)
			}
			u.turbineCLI = mock

			err := u.Execute(context.Background())
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
