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
	"reflect"
	"strings"
	"testing"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

func getConnectors(pipelineID int) []*meroxa.Connector {
	var connectors []*meroxa.Connector
	c := utils.GenerateConnector(pipelineID, "")
	connectors = append(connectors, &c)
	return connectors
}

func TestListConnectorsFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{name: "pipeline", required: false, shorthand: ""},
	}

	c := builder.BuildCobraCommand(&List{})

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

func TestListConnectorsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockListConnectorsClient(ctrl)
	logger := log.NewTestLogger()

	connectors := getConnectors(1)

	client.
		EXPECT().
		ListConnectors(ctx).
		Return(connectors, nil)

	lc := &List{
		client: client,
		logger: logger,
	}

	err := lc.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := utils.ConnectorsTable(connectors)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotConnectors []meroxa.Connector
	err = json.Unmarshal([]byte(gotJSONOutput), &gotConnectors)

	var cc []meroxa.Connector

	for _, c := range connectors {
		cc = append(cc, *c)
	}

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotConnectors, cc) {
		t.Fatalf("expected \"%v\", got \"%v\"", connectors, gotConnectors)
	}
}
