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

package pipelines

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/utils"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestCreatePipelineArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires a pipeline name"), name: ""},
		{args: []string{"pipeline-name"}, err: nil, name: "pipeline-name"},
	}

	for _, tt := range tests {
		cc := &Create{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestCreatePipelineFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "metadata", required: false, shorthand: "m", hidden: false},
		{name: "env", required: false, hidden: true},
	}

	c := builder.BuildCobraCommand(&Create{})

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

func TestCreatePipelineWithoutEnvironmentExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	pName := "my-pipeline"

	pi := &meroxa.CreatePipelineInput{
		Name: pName,
	}

	p := &meroxa.Pipeline{
		ID:    1,
		Name:  pName,
		State: "healthy",
	}

	p.Name = pName

	client.
		EXPECT().
		CreatePipeline(
			ctx,
			pi,
		).
		Return(p, nil)

	c := &Create{
		client: client,
		logger: logger,
	}

	c.args.Name = pi.Name

	err := c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating pipeline %q in %q environment...
Pipeline %q successfully created!
`, pName, string(meroxa.EnvironmentTypeCommon), pName)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotPipeline meroxa.Pipeline
	err = json.Unmarshal([]byte(gotJSONOutput), &gotPipeline)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotPipeline, *p) {
		t.Fatalf("expected \"%v\", got \"%v\"", *p, gotPipeline)
	}
}

func TestCreatePipelineWithEnvironmentExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()
	pName := "my-pipeline"
	env := "my-env"

	pi := &meroxa.CreatePipelineInput{
		Name:        pName,
		Environment: &meroxa.PipelineEnvironment{Name: env},
	}

	p := &meroxa.Pipeline{
		ID:   1,
		Name: pName,
		Environment: &meroxa.PipelineEnvironment{
			UUID: "2560fbcc-b9ee-461a-a959-fa5656422dc2",
			Name: env,
		},
		State: "healthy",
	}

	p.Name = pName

	client.
		EXPECT().
		CreatePipeline(
			ctx,
			pi,
		).
		Return(p, nil)

	c := &Create{
		client: client,
		logger: logger,
	}

	c.args.Name = pi.Name
	c.flags.Environment = pi.Environment.Name

	err := c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating pipeline %q in %q environment...
Pipeline %q successfully created!
`, pName, env, pName)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotPipeline meroxa.Pipeline
	err = json.Unmarshal([]byte(gotJSONOutput), &gotPipeline)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotPipeline, *p) {
		t.Fatalf("expected \"%v\", got \"%v\"", *p, gotPipeline)
	}
}
