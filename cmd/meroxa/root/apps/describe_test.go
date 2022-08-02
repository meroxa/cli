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
	"reflect"
	"strings"
	"testing"

	"github.com/volatiletech/null/v8"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestDescribeApplicationArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires app name or UUID"), name: ""},
		{args: []string{"ApplicationName"}, err: nil, name: "ApplicationName"},
	}

	for _, tt := range tests {
		ar := &Describe{}
		err := ar.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != ar.args.NameOrUUID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ar.args.NameOrUUID)
		}
	}
}

func TestDescribeApplicationExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	appName := "my-env"

	a := utils.GenerateApplication("")
	a.Name = appName

	a.Resources = []meroxa.ApplicationResource{
		{
			EntityIdentifier: meroxa.EntityIdentifier{
				Name: null.StringFrom("res1"),
			},
			Collection: meroxa.ResourceCollection{
				Name:   null.StringFrom("res1"),
				Source: null.StringFrom("source"),
			},
		},
		{
			EntityIdentifier: meroxa.EntityIdentifier{
				Name: null.StringFrom("res2"),
			},
			Collection: meroxa.ResourceCollection{
				Name:        null.StringFrom("res2"),
				Destination: null.StringFrom("destination"),
			},
		},
	}
	resources := []*meroxa.Resource{
		{Name: "res1", UUID: "abc-def", Type: meroxa.ResourceTypePostgres},
		{Name: "res2", UUID: "abc-def", Type: meroxa.ResourceTypeBigquery},
	}

	a.Connectors = []meroxa.EntityIdentifier{
		{Name: null.StringFrom("conn1")},
		{Name: null.StringFrom("conn2")},
	}
	connectors := []*meroxa.Connector{
		{Name: "conn1", ResourceName: "res1", Type: meroxa.ConnectorTypeSource, State: meroxa.ConnectorStateRunning},
		{Name: "conn2", ResourceName: "res2", Type: meroxa.ConnectorTypeDestination, State: meroxa.ConnectorStateRunning},
	}

	functions := []*meroxa.Function{
		{Name: "fun1", UUID: "abc-def", Status: meroxa.FunctionStatus{State: "running"}},
	}
	a.Functions = []meroxa.EntityIdentifier{
		{Name: null.StringFrom("fun1")},
	}

	client.EXPECT().GetApplication(ctx, a.Name).Return(&a, nil)
	client.EXPECT().GetResourceByNameOrID(ctx, "res1").Return(resources[0], nil)
	client.EXPECT().GetResourceByNameOrID(ctx, "res2").Return(resources[1], nil)
	client.EXPECT().GetConnectorByNameOrID(ctx, "conn1").Return(connectors[0], nil)
	client.EXPECT().GetConnectorByNameOrID(ctx, "conn2").Return(connectors[1], nil)
	client.EXPECT().GetFunction(ctx, "fun1").Return(functions[0], nil)

	dc := &Describe{
		client: client,
		logger: logger,
	}
	dc.args.NameOrUUID = a.Name

	err := dc.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := utils.AppTable(&a, resources, connectors, functions)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotApp meroxa.Application
	err = json.Unmarshal([]byte(gotJSONOutput), &gotApp)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotApp, a) {
		t.Fatalf("expected \"%v\", got \"%v\"", a, gotApp)
	}
}
