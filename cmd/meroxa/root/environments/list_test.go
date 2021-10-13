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
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

func TestListEnvironmentsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockListEnvironmentsClient(ctrl)
	logger := log.NewTestLogger()

	ee := &meroxa.Environment{
		Type:     "dedicated",
		Name:     "environment-1234",
		Provider: "aws",
		Region:   "aws:us-east",
		Status:   meroxa.EnvironmentStatus{State: "provisioned"},
		UUID:     "531428f7-4e86-4094-8514-d397d49026f7",
	}

	environments := []*meroxa.Environment{ee}

	client.
		EXPECT().
		ListEnvironments(ctx).
		Return(environments, nil)

	l := &List{
		client: client,
		logger: logger,
	}

	err := l.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := utils.EnvironmentsTable(environments, false)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotEnvironments []meroxa.Environment
	err = json.Unmarshal([]byte(gotJSONOutput), &gotEnvironments)

	var lp []meroxa.Environment

	for _, p := range environments {
		lp = append(lp, *p)
	}

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotEnvironments, lp) {
		t.Fatalf("expected \"%v\", got \"%v\"", environments, gotEnvironments)
	}
}
