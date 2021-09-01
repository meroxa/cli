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

package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/meroxa/cli/log"

	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/meroxa-go"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/utils"
)

func TestCreateConnectorArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: nil, name: ""},
		{args: []string{"conName"}, err: nil, name: "conName"},
	}

	for _, tt := range tests {
		cc := &Create{}
		err := cc.ParseArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestCreateConnectorFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "input", required: false},
		{name: "config", required: false, shorthand: "c"},
		{name: "from", required: false},
		{name: "to", required: false},
		{name: "metadata", required: false, shorthand: "m", hidden: true},
		{name: "pipeline", required: true},
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

func TestCreateConnectorExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockCreateConnectorClient(ctrl)
	logger := log.NewTestLogger()

	sourceName := "my-resource"

	c := &Create{
		client: client,
		logger: logger,
	}

	c.flags.Input = "foo"
	c.flags.Config = `{"key":"value"}`
	c.flags.Metadata = `{"metakey":"metavalue"}`
	c.flags.Source = sourceName
	c.flags.Pipeline = "my-pipeline"

	cr := utils.GenerateConnector(0, "")

	client.
		EXPECT().
		GetResourceByName(
			ctx,
			sourceName,
		).
		Return(&meroxa.Resource{ID: 123}, nil)

	client.
		EXPECT().
		CreateConnector(
			ctx,
			meroxa.CreateConnectorInput{
				Name:         "",
				ResourceID:   123,
				PipelineName: "my-pipeline",
				Configuration: map[string]interface{}{
					"key":   "value",
					"input": "foo",
				},
				Metadata: map[string]interface{}{
					"metakey":          "metavalue",
					"mx:connectorType": "source",
				},
			},
		).
		Return(&cr, nil)

	err := c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating connector from source %q...
Connector %q successfully created!
`, sourceName, cr.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotConnector meroxa.Connector
	err = json.Unmarshal([]byte(gotJSONOutput), &gotConnector)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotConnector, cr) {
		t.Fatalf("expected \"%v\", got \"%v\"", cr, gotConnector)
	}
}

func TestCreateConnectorExecutionWithPipelineID(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockCreateConnectorClient(ctrl)
	logger := log.NewTestLogger()

	sourceName := "my-resource"

	c := &Create{
		client: client,
		logger: logger,
	}

	c.flags.Input = "foo"
	c.flags.Config = `{"key":"value"}`
	c.flags.Metadata = `{"metakey":"metavalue"}`
	c.flags.Source = sourceName
	c.flags.Pipeline = "456"

	cr := utils.GenerateConnector(0, "")

	client.
		EXPECT().
		GetResourceByName(
			ctx,
			sourceName,
		).
		Return(&meroxa.Resource{ID: 123}, nil)

	client.
		EXPECT().
		CreateConnector(
			ctx,
			meroxa.CreateConnectorInput{
				Name:       "",
				ResourceID: 123,
				PipelineID: 456,
				Configuration: map[string]interface{}{
					"key":   "value",
					"input": "foo",
				},
				Metadata: map[string]interface{}{
					"metakey":          "metavalue",
					"mx:connectorType": "source",
				},
			},
		).
		Return(&cr, nil)

	err := c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating connector from source %q...
Connector %q successfully created!
`, sourceName, cr.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotConnector meroxa.Connector
	err = json.Unmarshal([]byte(gotJSONOutput), &gotConnector)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotConnector, cr) {
		t.Fatalf("expected \"%v\", got \"%v\"", cr, gotConnector)
	}
}

func TestCreateConnectorExecutionNoFlags(t *testing.T) {
	ctx := context.Background()
	c := &Create{}

	err := c.Execute(ctx)

	expected := "requires either a source (--from) or a destination (--to)"

	if err != nil && err.Error() != expected {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}
