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

package environments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestUpdateEnvironmentArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires environment name"), name: ""},
		{args: []string{"environment-name"}, err: nil, name: "environment-name"},
	}

	for _, tt := range tests {
		cc := &Update{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.NameOrUUID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.NameOrUUID)
		}
	}
}

func TestUpdateEnvironmentExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	r := &Update{
		client: client,
		logger: logger,
	}

	newName := "new-name"
	newConfig := []string{"a=b", "c=d"}
	e := utils.GenerateEnvironment("")
	r.args.NameOrUUID = e.Name
	r.flags.Name = newName
	r.flags.Config = newConfig
	input := &meroxa.UpdateEnvironmentInput{Name: newName, Configuration: stringSliceToMap(newConfig)}

	client.
		EXPECT().
		UpdateEnvironment(ctx, e.Name, input).
		Return(&e, nil)

	err := r.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(
		"Updating environment...\nPreflight checks have passed. Environment %q is being updated. Run `meroxa env describe %s` for status",
		e.Name,
		e.Name)

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
