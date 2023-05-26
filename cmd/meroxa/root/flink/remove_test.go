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

package flink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestRemoveFlinkJobArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires Flink Job name or UUID"), name: ""},
		{args: []string{"flink-job-name"}, err: nil, name: "flink-job-name"},
	}

	for _, tt := range tests {
		cc := &Remove{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.NameOrUUID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.NameOrUUID)
		}
	}
}

func TestRemoveFlinkJobExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	r := &Remove{
		client: client,
		logger: logger,
	}

	fj := utils.GenerateFlinkJob()
	r.args.NameOrUUID = fj.Name

	client.
		EXPECT().
		GetFlinkJob(ctx, r.args.NameOrUUID).
		Return(&fj, nil)

	client.
		EXPECT().
		DeleteFlinkJob(ctx, r.args.NameOrUUID).
		Return(nil)

	err := r.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Removing Flink Job %q...
Flink Job %q successfully removed
`, fj.Name, fj.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotFJ meroxa.FlinkJob
	err = json.Unmarshal([]byte(gotJSONOutput), &gotFJ)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotFJ, fj) {
		t.Fatalf("expected \"%v\", got \"%v\"", fj, gotFJ)
	}
}
