/*
Copyright © 2022 Meroxa Inc

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

package apps

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestListAppsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	aa := &meroxa.Application{
		Name:     "my-app",
		UUID:     "531428f7-4e86-4094-8514-d397d49026f7",
		Language: turbine.GoLang,
		Status:   meroxa.ApplicationStatus{State: meroxa.ApplicationStateRunning},
	}

	apps := []*meroxa.Application{aa}

	client.
		EXPECT().
		AddHeader("Meroxa-CLI-App", "1").
		Times(1)
	client.
		EXPECT().
		ListApplications(ctx).
		Return(apps, nil)

	l := &List{
		client: client,
		logger: logger,
	}

	err := l.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := display.AppsTable(apps, false)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotApps []*meroxa.Application
	err = json.Unmarshal([]byte(gotJSONOutput), &gotApps)

	var lp []*meroxa.Application

	lp = append(lp, apps...)

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotApps, lp) {
		t.Fatalf("expected \"%v\", got \"%v\"", apps, gotApps)
	}
}
