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

package apps

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/volatiletech/null/v8"

	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
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
	res1 := &http.Response{
		Body: io.NopCloser(strings.NewReader(log)),
	}
	res2 := &http.Response{
		Body: io.NopCloser(strings.NewReader(log)),
	}
	res3 := &http.Response{
		Body: io.NopCloser(strings.NewReader(log)),
	}

	a := utils.GenerateApplication("")
	a.Name = appName
	a.Resources = []meroxa.ApplicationResource{
		{meroxa.EntityIdentifier{Name: null.StringFrom("res1")}, //nolint:govet
			meroxa.ResourceCollectionView{Name: "res1", Source: "true"}},
		{meroxa.EntityIdentifier{Name: null.StringFrom("res2")}, //nolint:govet
			meroxa.ResourceCollectionView{Name: "res1", Destination: "true"}},
	}

	a.Connectors = []meroxa.EntityIdentifier{
		{Name: null.StringFrom("conn1")},
		{Name: null.StringFrom("conn2")},
	}
	connectors := []*utils.ExtendedConnector{
		{Connector: &meroxa.Connector{
			Name: "conn1", ResourceName: "res1", Type: meroxa.ConnectorTypeSource, State: meroxa.ConnectorStateRunning},
			Logs: log},
		{Connector: &meroxa.Connector{
			Name: "conn2", ResourceName: "res2", Type: meroxa.ConnectorTypeDestination, State: meroxa.ConnectorStateRunning},
			Logs: log},
	}

	functions := []*meroxa.Function{
		{Name: "fun1", UUID: "abc-def", Status: meroxa.FunctionStatus{State: meroxa.FunctionStateRunning}, Logs: log},
	}
	a.Functions = []meroxa.EntityIdentifier{
		{Name: null.StringFrom("fun1")},
	}

	client.EXPECT().GetApplication(ctx, a.Name).Return(&a, nil)
	client.EXPECT().GetConnectorByNameOrID(ctx, "conn1").
		Return(&meroxa.Connector{Name: "conn1", ResourceName: "res1"}, nil)
	client.EXPECT().GetConnectorLogs(ctx, "conn1").Return(res1, nil)
	client.EXPECT().GetConnectorByNameOrID(ctx, "conn2").
		Return(&meroxa.Connector{Name: "conn2", ResourceName: "res2"}, nil)
	client.EXPECT().GetConnectorLogs(ctx, "conn2").Return(res2, nil)
	client.EXPECT().GetFunction(ctx, "fun1").Return(functions[0], nil)
	client.EXPECT().GetFunctionLogs(ctx, "fun1").Return(res3, nil)

	dc := &Logs{
		client: client,
		logger: logger,
	}
	dc.args.NameOrUUID = a.Name

	err := dc.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := utils.AppLogsTable(a.Resources, connectors, functions)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf(cmp.Diff(wantLeveledOutput, gotLeveledOutput))
	}

	gotJSONOutput := logger.JSONOutput()
	var gotApp meroxa.Application
	err = json.Unmarshal([]byte(gotJSONOutput), &gotApp)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotApp, a) {
		t.Fatalf(cmp.Diff(a, gotApp))
	}
}
