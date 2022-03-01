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

	"github.com/meroxa/meroxa-go/pkg/meroxa"

	"github.com/meroxa/cli/utils"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestRemoveConnectorArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires connector name"), name: ""},
		{args: []string{"conName"}, err: nil, name: "conName"},
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

func TestRemoveConnectorExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	r := &Remove{
		client: client,
		logger: logger,
	}

	c := utils.GenerateConnector(0, "")
	r.args.Name = c.Name

	client.
		EXPECT().
		GetConnectorByNameOrID(ctx, c.Name).
		Return(&c, nil)

	client.
		EXPECT().
		DeleteConnector(ctx, c.Name).
		Return(nil)

	err := r.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Removing connector %q...
Connector %q successfully removed
`, c.Name, c.Name)

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
