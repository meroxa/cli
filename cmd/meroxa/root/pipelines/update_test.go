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

package pipelines

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

func TestUpdatePipelineArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires pipeline name"), name: ""},
		{args: []string{"my-pipeline"}, err: nil, name: "my-pipeline"},
	}

	for _, tt := range tests {
		cc := &Update{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestUpdatePipelineFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "state", required: false, shorthand: "", hidden: false},
		{name: "name", required: false, shorthand: "", hidden: false},
		{name: "metadata", required: false, shorthand: "m", hidden: false},
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

func TestUpdatePipelineExecutionNoFlags(t *testing.T) {
	ctx := context.Background()

	u := &Update{}

	err := u.Execute(ctx)

	expected := "requires either --name, --state or --metadata"

	if err != nil && err.Error() != expected {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestUpdatePipelineExecutionWithNewState(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	p := utils.GeneratePipeline()
	newState := meroxa.Action("pause")

	client.
		EXPECT().
		GetPipelineByName(ctx, p.Name).
		Return(&p, nil)

	client.
		EXPECT().
		UpdatePipelineStatus(ctx, p.ID, newState).
		Return(&p, nil)

	u := &Update{
		client: client,
		logger: logger,
	}

	u.args.Name = p.Name
	u.flags.State = string(newState)

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating pipeline %q...
Pipeline %q successfully updated!
`, u.args.Name, u.args.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotPipeline meroxa.Pipeline
	err = json.Unmarshal([]byte(gotJSONOutput), &gotPipeline)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotPipeline, p) {
		t.Fatalf("expected \"%v\", got \"%v\"", p, gotPipeline)
	}
}

func TestUpdatePipelineExecutionWithNewName(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	p := utils.GeneratePipeline()
	newName := "new-pipeline-name"
	pi := &meroxa.UpdatePipelineInput{
		Name: newName,
	}

	client.
		EXPECT().
		GetPipelineByName(ctx, p.Name).
		Return(&p, nil)

	client.
		EXPECT().
		UpdatePipeline(ctx, p.ID, pi).
		Return(&p, nil)

	u := &Update{
		client: client,
		logger: logger,
	}

	u.args.Name = p.Name
	u.flags.Name = newName

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating pipeline %q...
Pipeline %q successfully updated!
`, u.args.Name, u.args.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotPipeline meroxa.Pipeline
	err = json.Unmarshal([]byte(gotJSONOutput), &gotPipeline)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotPipeline, p) {
		t.Fatalf("expected \"%v\", got \"%v\"", p, gotPipeline)
	}
}

func TestUpdatePipelineExecutionWithNewMetadata(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	p := utils.GeneratePipeline()

	pi := &meroxa.UpdatePipelineInput{
		Metadata: map[string]interface{}{"key": "value"},
	}

	client.
		EXPECT().
		GetPipelineByName(ctx, p.Name).
		Return(&p, nil)

	client.
		EXPECT().
		UpdatePipeline(ctx, p.ID, pi).
		Return(&p, nil)

	u := &Update{
		client: client,
		logger: logger,
	}

	u.args.Name = p.Name
	u.flags.Metadata = "{\"key\": \"value\"}"

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating pipeline %q...
Pipeline %q successfully updated!
`, u.args.Name, u.args.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotPipeline meroxa.Pipeline
	err = json.Unmarshal([]byte(gotJSONOutput), &gotPipeline)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotPipeline, p) {
		t.Fatalf("expected \"%v\", got \"%v\"", p, gotPipeline)
	}
}
