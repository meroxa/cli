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
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
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

	c := &Logs{
		client: client,
		logger: logger,
	}

	c.args.UUID = buildUUID

	var responseDetails = io.NopCloser(bytes.NewReader([]byte(
		`[2021-04-29T12:16:42Z] Beep boop, robots doing build things`,
	)))

	var httpResponse = &http.Response{
		StatusCode: 200,
		Body:       responseDetails,
	}

	client.
		EXPECT().
		GetBuildLogs(ctx, buildUUID).
		Return(httpResponse, nil)

	err := c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := "[2021-04-29T12:16:42Z] Beep boop, robots doing build things"

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}
}
