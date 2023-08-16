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

package builds

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestLogsBuildArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires build UUID"), name: ""},
		{args: []string{"buildUUID"}, err: nil, name: "buildUUID"},
	}

	for _, tt := range tests {
		cc := &Logs{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.UUID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.UUID)
		}
	}
}

func TestLogsBuildExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	buildUUID := "236d6e81-6a22-4805-b64f-3fa0a57fdbdc"

	l := &Logs{
		client: client,
		logger: logger,
	}

	l.args.UUID = buildUUID

	buildLog := &meroxa.Logs{
		Data: []meroxa.LogData{
			{
				Timestamp: time.Now().UTC(),
				Log:       "Beep boop, robots doing build things",
				Source:    "function build",
			},
		},
		Metadata: meroxa.Metadata{
			End:   time.Now().UTC(),
			Start: time.Now().UTC().Add(-12 * time.Hour),
			Limit: 10,
		},
	}

	client.
		EXPECT().
		GetBuildLogsV2(ctx, buildUUID).
		Return(buildLog, nil)

	err := l.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := display.BuildLogsTableV2(buildLog)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf(cmp.Diff(wantLeveledOutput, gotLeveledOutput))
	}

	gotJSONOutput := logger.JSONOutput()
	var gotBuildLog meroxa.Logs
	err = json.Unmarshal([]byte(gotJSONOutput), &gotBuildLog)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotBuildLog, *buildLog) {
		t.Fatalf(cmp.Diff(*buildLog, gotBuildLog))
	}
}
