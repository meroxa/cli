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
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	turbineMock "github.com/meroxa/cli/cmd/meroxa/turbine/mock"

	"github.com/meroxa/turbine-core/pkg/ir"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestRemoveAppArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires application name"), name: ""},
		{args: []string{"application-name"}, err: nil, name: "application-name"},
	}

	for _, tt := range tests {
		cc := &Remove{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.NameOrUUID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.NameOrUUID)
		}
	}
}

func TestRemoveAppExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	r := &Remove{
		client: client,
		logger: logger,
	}

	app := utils.GenerateApplication("")
	r.args.NameOrUUID = app.Name
	r.flags.Force = true

	res := &http.Response{
		StatusCode: http.StatusNoContent,
	}

	client.
		EXPECT().
		DeleteApplicationEntities(ctx, r.args.NameOrUUID).
		Return(res, nil)

	err := r.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Removing application %q...
Application %q successfully removed
`, app.Name, app.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()

	var gotResponse *http.Response
	err = json.Unmarshal([]byte(gotJSONOutput), &gotResponse)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotResponse, res) {
		t.Fatalf("expected \"%v\", got \"%v\"", gotResponse, res)
	}

	r.flags.Path = "does not matter"
	err = r.Execute(ctx)
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestRemoveAppExecutionWithPath(t *testing.T) {
	os.Setenv("UNIT_TEST", "1")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	mockTurbineCLI := turbineMock.NewMockCLI(ctrl)

	app := utils.GenerateApplication("")

	i := &Init{
		logger: logger,
	}
	path := filepath.Join(os.TempDir(), uuid.NewString())
	i.args.appName = app.Name
	i.flags.Path = path
	i.flags.Lang = string(ir.GoLang)
	i.flags.SkipModInit = true
	i.flags.ModVendor = false
	err := i.Execute(ctx)
	defer func(path string) {
		os.RemoveAll(path)
	}(path)
	require.NoError(t, err)

	logger = log.NewTestLogger()
	r := &Remove{
		client:     client,
		logger:     logger,
		turbineCLI: mockTurbineCLI,
	}

	r.flags.Path = filepath.Join(path, app.Name)
	r.flags.Force = true

	res := &http.Response{
		StatusCode: http.StatusNoContent,
	}

	mockTurbineCLI.EXPECT().GetVersion(ctx).Return("1.0", nil)
	client.EXPECT().AddHeader("Meroxa-CLI-App-Lang", string(ir.GoLang)).Times(1)
	client.EXPECT().AddHeader("Meroxa-CLI-App-Version", gomock.Any()).Times(1)
	client.
		EXPECT().
		DeleteApplicationEntities(ctx, app.Name).
		Return(res, nil)

	err = r.Execute(ctx)
	os.Setenv("UNIT_TEST", "")

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Removing application %q...
Application %q successfully removed
`, app.Name, app.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()

	var gotResponse *http.Response
	err = json.Unmarshal([]byte(gotJSONOutput), &gotResponse)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotResponse, res) {
		t.Fatalf("expected \"%v\", got \"%v\"", gotResponse, res)
	}
}
