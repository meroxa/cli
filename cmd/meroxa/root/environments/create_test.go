/*
Copyright © 2021 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package environments

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestCreateEnvironmentArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: nil, name: ""},
		{args: []string{"env-name"}, err: nil, name: "env-name"},
	}

	for _, tt := range tests {
		cc := &Create{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestCreateEnvironmentFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "type", required: false, hidden: false},
		{name: "provider", required: false, hidden: false},
		{name: "region", required: false, hidden: false},
		{name: "config", shorthand: "c", required: false, hidden: false},
	}

	c := builder.BuildCobraCommand(&Create{})

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		} else {
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
}

func TestCreateEnvironmentExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	c := &Create{
		client: client,
		logger: logger,
	}

	c.args.Name = "my-env"
	c.flags.Type = "dedicated"
	c.flags.Provider = "aws"
	c.flags.Region = "aws"
	c.flags.Config = []string{"aws_access_key_id=my_access_key", "aws_access_secret=my_access_secret"}

	cfg := stringSliceToMap(c.flags.Config)

	e := &meroxa.CreateEnvironmentInput{
		Type:          meroxa.EnvironmentType(c.flags.Type),
		Provider:      meroxa.EnvironmentProvider(c.flags.Provider),
		Name:          c.args.Name,
		Configuration: cfg,
		Region:        meroxa.EnvironmentRegion(c.flags.Region),
	}

	rE := &meroxa.Environment{
		UUID:          "602c4608-0c71-43c7-9d0a-0cb2ab9c9ccd",
		Name:          e.Name,
		Provider:      e.Provider,
		Region:        e.Region,
		Type:          e.Type,
		Configuration: e.Configuration,
		Status: meroxa.EnvironmentViewStatus{
			State: meroxa.EnvironmentStateProvisioning,
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	client.
		EXPECT().
		CreateEnvironment(
			ctx,
			e,
		).
		Return(rE, nil)

	err := c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf("Provisioning environment...\n"+
		"Environment %q is being provisioned. Run `meroxa env describe %s` for status", e.Name, e.Name)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotEnviroment meroxa.Environment
	err = json.Unmarshal([]byte(gotJSONOutput), &gotEnviroment)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotEnviroment, *rE) {
		t.Fatalf("expected \"%v\", got \"%v\"", *rE, gotEnviroment)
	}
}
