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

package api

import (
	"bytes"
	"context"
	"errors"
<<<<<<< HEAD
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
=======
	"io/ioutil"
	"net/http"
>>>>>>> 451b598 (refactor & feature: Convert API to builder command)
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"
)

func TestDescribeAPIArgs(t *testing.T) {
	tests := []struct {
		args   []string
		err    error
		method string
		path   string
		body   string
	}{
		{
			args: nil,
			err:  errors.New("requires METHOD and PATH"),
		},
		{
			args:   []string{"GET", "/v1/endpoints"},
			err:    nil,
			method: "GET",
			path:   "/v1/endpoints",
		},
		{
			args:   []string{"get", "/v1/endpoints"}, // lowercase
			err:    nil,
			method: "GET",
			path:   "/v1/endpoints",
		},
		{
			args: []string{
				"POST",
				"/v1/endpoints",
				"'{\"protocol\": \"HTTP\", \"stream\": \"resource-2-499379-public.accounts\", \"name\": \"1234\"}'"},
			err:    nil,
			method: "POST",
			path:   "/v1/endpoints",
			body:   "'{\"protocol\": \"HTTP\", \"stream\": \"resource-2-499379-public.accounts\", \"name\": \"1234\"}'",
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
	client := mock.NewMockAPIClient(ctrl)
	logger := log.NewTestLogger()

	a := &API{
		client: client,
		logger: logger,
	}
	a.args.Method = "GET"
	a.args.Path = "/v1/my-path"

<<<<<<< HEAD
	bodyResponse := `{ "key": "value" }`

	var httpResponse = &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(bodyResponse))),
=======
	var httpResponse = &http.Response{
		StatusCode: 200,
		Body: ioutil.NopCloser(bytes.NewReader([]byte(
			`{ "key": "value" }`,
		))),
>>>>>>> 451b598 (refactor & feature: Convert API to builder command)
	}

	client.
		EXPECT().
		MakeRequest(
			ctx,
			a.args.Method,
			a.args.Path,
			"",
			nil,
		).
		Return(httpResponse, nil)

	err := a.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

<<<<<<< HEAD
	expectedBody := `{
	"key": "value"
}`

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`> %s %s
< %s 
%s
`, a.args.Method, a.args.Path, httpResponse.Status, expectedBody)
=======
	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := `> GET /v1/my-path
<  
{
	"key": "value"
}
`
>>>>>>> 451b598 (refactor & feature: Convert API to builder command)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}
<<<<<<< HEAD

	gotJSONOutput := strings.TrimSpace(logger.JSONOutput())

	if !strings.Contains(expectedBody, gotJSONOutput) {
		t.Fatalf("expected \"%v\", got \"%v\"", expectedBody, gotJSONOutput)
	}
=======
>>>>>>> 451b598 (refactor & feature: Convert API to builder command)
}
