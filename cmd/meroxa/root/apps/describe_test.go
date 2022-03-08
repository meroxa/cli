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
	"reflect"
	"strings"
	"testing"

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
		{args: nil, err: errors.New("requires app name"), name: ""},
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

	a := utils.GenerateApplication()
	a.Name = appName

	client.
		EXPECT().
		GetApplication(
			ctx,
			a.Name,
		).
		Return(&a, nil)

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
	wantLeveledOutput := utils.AppTable(&a)

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

func TestDescribeApplicationExecutionWithFunctions(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	appName := "my-app-with-funcs"

	a := utils.GenerateApplication()
	a.Name = appName
	a.Functions = []meroxa.FunctionIdentifier{
		{Name: "fun1"},
		{Name: "fun2"},
	}

	client.
		EXPECT().
		GetApplication(
			ctx,
			a.Name,
		).
		Return(&a, nil)

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
	wantLeveledOutput := utils.AppTable(&a)

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
