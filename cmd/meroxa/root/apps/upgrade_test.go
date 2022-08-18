package apps

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/root/nop"
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
		lang   string
		vendor bool
		err    error
	}{
		{
			desc:   "Successful upgrade without vendor",
			vendor: false,
			err:    nil,
		},
		{
			desc:   "Successful upgrade with vendor",
			vendor: true,
			err:    nil,
		},
		{
			desc:   "Unsuccessful upgrade",
			lang:   GoLang,
			vendor: true,
			err:    fmt.Errorf("not good"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)

			u := &Upgrade{}
			u.Logger(log.NewTestLogger())
			u.flags.Path = "/does/not/matter"
			u.run = &nop.Nop{}

			mock := mockturbinecli.NewMockCLI(mockCtrl)
			if tt.err == nil {
				mock.EXPECT().Upgrade(tt.vendor)
			} else {
				mock.EXPECT().Upgrade(tt.vendor).Return(tt.err)
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
