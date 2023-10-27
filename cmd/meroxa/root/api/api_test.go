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

package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	basicMock "github.com/meroxa/cli/cmd/meroxa/global/mock"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
)

func TestDescribeAPIArgs(t *testing.T) {
	tests := []struct {
		args   []string
		err    error
		method string
		path   string
		body   interface{}
	}{
		{
			args:   []string{"GET", "/v1/resources"},
			err:    nil,
			method: "GET",
			path:   "/v1/resources",
			body:   nil,
		},
		{
			args:   []string{"get", "/v1/resources"}, // lowercase
			err:    nil,
			method: "GET",
			path:   "/v1/resources",
			body:   nil,
		},
		{
			args: []string{
				"POST",
				"/v1/resources",
				`'{"type":"postgres", "name":"pg", "url":"postgres://u:p@127.0.01:5432/db"}'`,
			},
			err:    nil,
			method: "POST",
			path:   "/v1/resources",
			body:   `'{"type":"postgres", "name":"pg", "url":"postgres://u:p@127.0.01:5432/db"}'`,
		},
		{
			args: nil,
			err:  errors.New("requires METHOD and PATH"),
		},
	}

	for _, tt := range tests {
		a := &API{}
		err := a.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.method != a.args.Method {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.method, a.args.Method)
		}

		if tt.path != a.args.Path {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.path, a.args.Path)
		}

		if tt.body != a.args.Body {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.body, a.args.Body)
		}
	}
}

func TestAPIExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := basicMock.NewMockBasicClient(ctrl)
	logger := log.NewTestLogger()

	a := &API{
		client: client,
		logger: logger,
	}
	a.args.Method = "GET"
	a.args.Path = "/api/collections/apps/records"
	a.args.ID = "04b0d690-dd44-4df3-8"
	a.args.Body = "somebody"

	bodyResponse := `{ "key": "value" }`

	httpResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Body:       io.NopCloser(bytes.NewReader([]byte(bodyResponse))),
	}

	client.EXPECT().URLRequest(ctx, "GET", "/api/collections/apps/records", "", nil, nil, nil).Return(
		httpResponse,
		nil,
	)

	err := a.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	expectedBody := `{
	"key": "value"
}`

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`> %s %s
< %s %s
%s
`, a.args.Method, a.args.Path, httpResponse.Status, httpResponse.Proto, expectedBody)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := strings.TrimSpace(logger.JSONOutput())

	if !strings.Contains(expectedBody, gotJSONOutput) {
		t.Fatalf("expected \"%v\", got \"%v\"", expectedBody, gotJSONOutput)
	}
}
