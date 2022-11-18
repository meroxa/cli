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
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/cli/utils/display"
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
				Name: "res1",
			},
			Collection: meroxa.ResourceCollection{
				Name:   "res1",
				Source: "source",
			},
		},
		{
			EntityIdentifier: meroxa.EntityIdentifier{
				Name: "res2",
			},
			Collection: meroxa.ResourceCollection{
				Name:        "res2",
				Destination: "destination",
			},
		},
		{
			EntityIdentifier: meroxa.EntityIdentifier{
				Name: "res3",
			},
			Collection: meroxa.ResourceCollection{
				Name:        "res3",
				Destination: "destination",
			},
		},
	}
	a.Connectors = []meroxa.EntityDetails{
		{EntityIdentifier: meroxa.EntityIdentifier{Name: "conn1"}},
		{EntityIdentifier: meroxa.EntityIdentifier{Name: "conn2"}},
		{EntityIdentifier: meroxa.EntityIdentifier{Name: "conn3"}},
	}
	a.Functions = []meroxa.EntityDetails{
		{EntityIdentifier: meroxa.EntityIdentifier{Name: "fun1"}},
	}

	client.EXPECT().GetApplication(ctx, a.Name).Return(&a, nil)

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
	wantLeveledOutput := display.AppTable(&a)

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

	dc.flags.Path = "does not matter"
	err = dc.Execute(ctx)
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestDescribeApplicationExecutionWithPath(t *testing.T) {
	os.Setenv("UNIT_TEST", "1")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	appName := "my-env"

	i := &Init{
		logger: logger,
	}
	path := "tmp" + uuid.NewString()
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

	a := utils.GenerateApplication("")
	a.Name = appName

	a.Resources = []meroxa.ApplicationResource{
		{
			EntityIdentifier: meroxa.EntityIdentifier{
				Name: "res1",
			},
			Collection: meroxa.ResourceCollection{
				Name:   "res1",
				Source: "source",
			},
		},
		{
			EntityIdentifier: meroxa.EntityIdentifier{
				Name: "res2",
			},
			Collection: meroxa.ResourceCollection{
				Name:        "res2",
				Destination: "destination",
			},
		},
		{
			EntityIdentifier: meroxa.EntityIdentifier{
				Name: "res3",
			},
			Collection: meroxa.ResourceCollection{
				Name:        "res3",
				Destination: "destination",
			},
		},
	}
	a.Connectors = []meroxa.EntityDetails{
		{EntityIdentifier: meroxa.EntityIdentifier{Name: "conn1"}},
		{EntityIdentifier: meroxa.EntityIdentifier{Name: "conn2"}},
		{EntityIdentifier: meroxa.EntityIdentifier{Name: "conn3"}},
	}
	a.Functions = []meroxa.EntityDetails{
		{EntityIdentifier: meroxa.EntityIdentifier{Name: "fun1"}},
	}

	client.EXPECT().AddHeader("Meroxa-CLI-App", "1").Times(1)
	client.EXPECT().AddHeader("Meroxa-CLI-App-Lang", turbine.GoLang).Times(1)
	client.EXPECT().AddHeader("Meroxa-CLI-App-Version", gomock.Any()).Times(1)
	client.EXPECT().GetApplication(ctx, a.Name).Return(&a, nil)

	dc := &Describe{
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
	wantLeveledOutput := display.AppTable(&a)

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
