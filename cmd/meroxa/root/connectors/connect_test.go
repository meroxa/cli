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

func TestConnectFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{name: "from", required: true, shorthand: ""},
		{name: "to", required: true, shorthand: ""},
		{name: "config", required: false, shorthand: "c"},
		{name: "input", required: false, shorthand: ""},
		{name: "pipeline", required: true, shorthand: ""},
	}

	c := builder.BuildCobraCommand(&Connect{})

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
	}
}

//nolint:funlen
func TestConnectExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	c := &Connect{
		client: client,
		logger: logger,
	}

	rSource := utils.GenerateResource()
	rDestination := utils.GenerateResource()

	c.flags.Input = "my-resource.Table"
	c.flags.Config = `{"key":"value"}`
	c.flags.Source = rSource.Name
	c.flags.Destination = rDestination.Name
	c.flags.Pipeline = "my-pipeline"

	cSource := utils.GenerateConnector("", "")
	cSource.Type = meroxa.ConnectorTypeSource

	cDestination := utils.GenerateConnector("", "")
	cDestination.Type = meroxa.ConnectorTypeDestination

	// Create source
	client.
		EXPECT().
		GetResourceByNameOrID(
			ctx,
			rSource.Name,
		).
		Return(&rSource, nil)

	client.
		EXPECT().
		CreateConnector(
			ctx,
			&meroxa.CreateConnectorInput{
				Name:         "",
				ResourceName: rSource.Name,
				Configuration: map[string]interface{}{
					"key": "value",
				},
				Metadata:     map[string]interface{}{},
				PipelineName: c.flags.Pipeline,
				Input:        "my-resource.Table",
				Type:         meroxa.ConnectorTypeSource,
			},
		).
		Return(&cSource, nil)

	// Create destination
	client.
		EXPECT().
		GetResourceByNameOrID(
			ctx,
			rDestination.Name,
		).
		AnyTimes().
		Return(&rDestination, nil)

	client.
		EXPECT().
		CreateConnector(
			ctx,
			&meroxa.CreateConnectorInput{
				Name:         "",
				ResourceName: rDestination.Name,
				Configuration: map[string]interface{}{
					"key": "value",
				},
				Metadata:     map[string]interface{}{},
				PipelineName: c.flags.Pipeline,
				Input:        "my-resource.Table",
				Type:         meroxa.ConnectorTypeDestination,
			},
		).
		Return(&cDestination, nil)

	err := c.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Creating connector from source %q in pipeline %q...
Creating connector to destination %q in pipeline %q...
Source connector %q and destination connector %q successfully created!
`, rSource.Name, c.flags.Pipeline, rDestination.Name, c.flags.Pipeline, cSource.Name, cDestination.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotConnectors []meroxa.Connector
	err = json.Unmarshal([]byte(gotJSONOutput), &gotConnectors)

	var connectors []meroxa.Connector
	connectors = append(connectors, cSource, cDestination)

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotConnectors, connectors) {
		t.Fatalf("expected \"%v\", got \"%v\"", connectors, gotConnectors)
	}
}

func TestConnectExecutionNoFlags(t *testing.T) {
	ctx := context.Background()
	c := &Connect{}

	err := c.Execute(ctx)

	expected := "requires either a source (--from) or a destination (--to)"

	if err != nil && err.Error() != expected {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}
