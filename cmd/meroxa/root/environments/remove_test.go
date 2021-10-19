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
	"errors"
	"strings"

	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

func TestRemoveEnvironmentArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires environment name"), name: ""},
		{args: []string{"environment-name"}, err: nil, name: "environment-name"},
	}

	for _, tt := range tests {
		cc := &Remove{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestRemoveEnvironmentExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockRemoveEnvironmentClient(ctrl)
	logger := log.NewTestLogger()

	r := &Remove{
		client: client,
		logger: logger,
	}

	e := utils.GenerateEnvironment("")
	r.args.Name = e.Name

	client.
		EXPECT().
		DeleteEnvironment(ctx, e.Name).
		Return(&e, nil)

	err := r.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf("Environment %q is being removed...\n", e.Name)
	wantLeveledOutput += fmt.Sprintf("Run `meroxa env describe %s` for status.", e.Name)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotEnvironment meroxa.Environment
	err = json.Unmarshal([]byte(gotJSONOutput), &gotEnvironment)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotEnvironment, e) {
		t.Fatalf("expected \"%v\", got \"%v\"", e, gotEnvironment)
	}
}
