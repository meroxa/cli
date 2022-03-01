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

package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestUpdateConnectorArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires connector name"), name: ""},
		{args: []string{"conName"}, err: nil, name: "conName"},
	}

	for _, tt := range tests {
		cc := &Update{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.NameOrID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.NameOrID)
		}
	}
}

func TestUpdateConnectorFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "config", required: false, shorthand: "c", hidden: false},
		{name: "name", required: false, shorthand: "", hidden: false},
		{name: "state", required: false, shorthand: "", hidden: false},
	}

	c := builder.BuildCobraCommand(&Update{})

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		}

		if f.shorthand != cf.Shorthand {
			t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
		}

		if f.required && !utils.IsFlagRequired(cf) {
			t.Fatalf("expected flag \"%s\" to be required", f.name)
		}

		if cf.Hidden != f.hidden {
			if cf.Hidden {
				t.Fatalf("expected flag \"%s\" not to be hidden", f.name)
			} else {
				t.Fatalf("expected flag \"%s\" to be hidden", f.name)
			}
		}
	}
}

func TestUpdateConnectorExecutionNoFlags(t *testing.T) {
	ctx := context.Background()
	u := &Update{}

	err := u.Execute(ctx)

	expected := "requires either --config, --name or --state"

	if err != nil && err.Error() != expected {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestUpdateConnectorExecutionWithNewState(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	u := &Update{
		client: client,
		logger: logger,
	}

	c := utils.GenerateConnector(0, "")
	u.args.NameOrID = c.Name
	u.flags.State = "pause"

	client.
		EXPECT().
		UpdateConnectorStatus(ctx, u.args.NameOrID, meroxa.Action(u.flags.State)).
		Return(&c, nil)

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating connector %q...
Connector %q successfully updated!
`, u.args.NameOrID, u.args.NameOrID)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotConnector meroxa.Connector
	err = json.Unmarshal([]byte(gotJSONOutput), &gotConnector)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotConnector, c) {
		t.Fatalf("expected \"%v\", got \"%v\"", c, gotConnector)
	}
}

func TestUpdateConnectorExecutionWithNewName(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	u := &Update{
		client: client,
		logger: logger,
	}

	c := utils.GenerateConnector(0, "")
	u.args.NameOrID = c.Name

	newName := "new-name"
	u.flags.Name = newName
	cu := meroxa.UpdateConnectorInput{
		Name: newName,
	}

	client.
		EXPECT().
		UpdateConnector(ctx, u.args.NameOrID, &cu).
		Return(&c, nil)

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating connector %q...
Connector %q successfully updated!
`, u.args.NameOrID, u.args.NameOrID)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotConnector meroxa.Connector
	err = json.Unmarshal([]byte(gotJSONOutput), &gotConnector)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotConnector, c) {
		t.Fatalf("expected \"%v\", got \"%v\"", c, gotConnector)
	}
}

func TestUpdateConnectorExecutionWithNewConfig(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	u := &Update{
		client: client,
		logger: logger,
	}

	c := utils.GenerateConnector(0, "")
	u.args.NameOrID = c.Name

	newConfig := "{\"table.name.format\":\"public.copy\"}"
	cfg := map[string]interface{}{}

	u.flags.Config = newConfig
	err := json.Unmarshal([]byte(u.flags.Config), &cfg)

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	cu := meroxa.UpdateConnectorInput{
		Configuration: cfg,
	}

	client.
		EXPECT().
		UpdateConnector(ctx, u.args.NameOrID, &cu).
		Return(&c, nil)

	err = u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating connector %q...
Connector %q successfully updated!
`, u.args.NameOrID, u.args.NameOrID)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotConnector meroxa.Connector
	err = json.Unmarshal([]byte(gotJSONOutput), &gotConnector)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotConnector, c) {
		t.Fatalf("expected \"%v\", got \"%v\"", c, gotConnector)
	}
}
