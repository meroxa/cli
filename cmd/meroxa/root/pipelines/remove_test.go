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

	"github.com/meroxa/meroxa-go"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
)

func TestRemovePipelineArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires pipeline name"), name: ""},
		{args: []string{"pipeline-name"}, err: nil, name: "pipeline-name"},
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

func TestRemovePipelineExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockRemovePipelineClient(ctrl)
	logger := log.NewTestLogger()

	r := &Remove{
		client: client,
		logger: logger,
	}

	p := utils.GeneratePipeline()
	r.args.Name = p.Name

	client.
		EXPECT().
		GetPipelineByName(ctx, p.Name).
		Return(&p, nil)

	client.
		EXPECT().
		DeletePipeline(ctx, p.ID).
		Return(nil)

	err := r.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Removing pipeline %q...
Pipeline %q successfully removed
`, r.args.Name, r.args.Name)

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
