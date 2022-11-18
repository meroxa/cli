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
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestApplicationLogsArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires app name or UUID"), name: ""},
		{args: []string{"ApplicationName"}, err: nil, name: "ApplicationName"},
	}

	for _, tt := range tests {
		ar := &Logs{}
		err := ar.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != ar.args.NameOrUUID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ar.args.NameOrUUID)
		}
	}
}

func TestApplicationLogsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	appName := "my-app-with-funcs"
	log := "hello world"

	appLogs := &meroxa.ApplicationLogs{
		ConnectorLogs:  map[string]string{"res1": log},
		FunctionLogs:   map[string]string{"fun1": log},
		DeploymentLogs: map[string]string{"uu-id": log},
	}

	client.EXPECT().GetApplicationLogs(ctx, appName).Return(appLogs, nil)

	dc := &Logs{
		client: client,
		logger: logger,
	}
	dc.args.NameOrUUID = appName

	err := dc.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := display.AppLogsTable(appLogs)

	// N.B. This comparison is undeterminstic when the test data map contains
	//      more than one key. Maps in golang are not guaranteed ordered so the result
	//      can be different.
	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf(cmp.Diff(wantLeveledOutput, gotLeveledOutput))
	}

	gotJSONOutput := logger.JSONOutput()
	var gotAppLogs meroxa.ApplicationLogs
	err = json.Unmarshal([]byte(gotJSONOutput), &gotAppLogs)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotAppLogs, *appLogs) {
		t.Fatalf(cmp.Diff(*appLogs, gotAppLogs))
	}

	dc.flags.Path = "does not matter"
	err = dc.Execute(ctx)
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestApplicationLogsExecutionWithPath(t *testing.T) {
	os.Setenv("UNIT_TEST", "1")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	appName := "my-app-with-funcs"
	log := "hello world"

	i := &Init{
		logger: logger,
	}
	path, _ := filepath.Abs("tmp" + uuid.NewString())
	i.args.appName = appName
	i.flags.Path = path
	i.flags.Lang = turbine.GoLang
	i.flags.SkipModInit = true
	i.flags.ModVendor = false
	err := i.Execute(ctx)
	defer func(path string) {
		os.RemoveAll(path)
	}(path)
	require.NoError(t, err)

	appLogs := &meroxa.ApplicationLogs{
		ConnectorLogs:  map[string]string{"res1": log},
		FunctionLogs:   map[string]string{"fun1": log},
		DeploymentLogs: map[string]string{"uu-id": log},
	}

	client.EXPECT().AddHeader("Meroxa-CLI-App-Lang", turbine.GoLang).Times(1)
	client.EXPECT().AddHeader("Meroxa-CLI-App-Version", gomock.Any()).Times(1)
	client.EXPECT().GetApplicationLogs(ctx, appName).Return(appLogs, nil)

	dc := &Logs{
		client: client,
		logger: logger,
	}
	dc.flags.Path = filepath.Join(path, appName)

	err = dc.Execute(ctx)
	os.Setenv("UNIT_TEST", "")
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := display.AppLogsTable(appLogs)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf(cmp.Diff(wantLeveledOutput, gotLeveledOutput))
	}

	gotJSONOutput := logger.JSONOutput()
	var gotAppLogs meroxa.ApplicationLogs
	err = json.Unmarshal([]byte(gotJSONOutput), &gotAppLogs)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotAppLogs, *appLogs) {
		t.Fatalf(cmp.Diff(*appLogs, gotAppLogs))
	}
}
