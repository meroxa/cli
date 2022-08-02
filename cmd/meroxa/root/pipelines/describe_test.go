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
	"reflect"
	"strings"
	"testing"

	"github.com/meroxa/cli/utils/display"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestDescribePipelineArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires pipeline name"), name: ""},
		{args: []string{"pipeline-name"}, err: nil, name: "pipeline-name"},
	}

	for _, tt := range tests {
		ar := &Describe{}
		err := ar.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != ar.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ar.args.Name)
		}
	}
}

func TestDescribePipelineExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	p := utils.GeneratePipelineWithEnvironment()
	client.
		EXPECT().
		GetPipelineByName(
			ctx,
			p.Name,
		).
		Return(&p, nil)

	dp := &Describe{
		client: client,
		logger: logger,
	}
	dp.args.Name = p.Name

	err := dp.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := display.PipelineTable(&p)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
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
