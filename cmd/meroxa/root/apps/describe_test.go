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

package apps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	basicMock "github.com/meroxa/cli/cmd/meroxa/global/mock"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/log"
)

const (
	body = `
	{
		"page": 1,
		"perPage": 30,
		"totalItems": 1,
		"totalPages": 1,
		"items": [
		  {
			"collectionId": "77byam8idl1rv8b",
			"collectionName": "conduitapps",
			"config": null,
			"created": "2024-04-01 20:13:20.111Z",
			"deployment_id": [
			  "tcsmsunmfo5v8kw"
			],
			"id": "lxjcdlsvet3aeoe",
			"name": "test-pipeline-1",
			"pipeline_enriched": "pipeline enriched settings",
			"pipeline_filename": "test-pipeline-1.yaml",
			"pipeline_original": "pipeline original settings",
			"state": "provisioned",
			"stream_provider": "kafka",
			"updated": "2024-04-01 20:13:20.111Z"
		  }
		]
	  }

	`
)

func TestDescribeApplicationArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires app name or UUID"), name: ""},
		{args: []string{"ApplicationName"}, err: nil, name: "ApplicationName"},
	}

	for _, tt := range tests {
		ar := &Describe{}
		err := ar.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != ar.args.nameOrUUID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ar.args.nameOrUUID)
		}
	}
}

func TestDescribeApplicationExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := basicMock.NewMockBasicClient(ctrl)
	logger := log.NewTestLogger()
	path := filepath.Join(os.TempDir(), uuid.NewString())
	appTime := AppTime{}
	err := appTime.UnmarshalJSON([]byte(`"2024-04-01 20:13:20.111Z"`))
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	a := &Application{
		Name:              "test-pipeline-1",
		State:             "provisioned",
		Created:           appTime,
		Updated:           appTime,
		PipelineFilenames: "test-pipeline-1.yaml",
		PipelineEnriched:  "pipeline enriched settings",
		PipelineOriginal:  "pipeline original settings",
		DeploymentID:      []string{"tcsmsunmfo5v8kw"},
		ApplicationSpec:   "kafka",
	}

	filter := &url.Values{}
	filter.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", a.Name, a.Name))

	httpResp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(body)),
		Status:     "200 OK",
		StatusCode: 200,
	}
	client.EXPECT().CollectionRequest(ctx, "GET", applicationCollection, "", nil, *filter).Return(
		httpResp,
		nil,
	)

	dc := &Describe{
		client: client,
		logger: logger,
		args:   struct{ nameOrUUID string }{nameOrUUID: a.Name},
		flags: struct {
			Path string "long:\"path\" usage:\"Path to the app directory (default is local directory)\""
		}{Path: filepath.Join(path, "appName")},
	}

	err = dc.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	var gotApp Application
	err = json.Unmarshal([]byte(logger.JSONOutput()), &gotApp)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if gotApp.Name != a.Name {
		t.Fatalf("expected \"%s\" got \"%s\"", a.Name, gotApp.Name)
	}
	if gotApp.State != a.State {
		t.Fatalf("expected \"%s\" got \"%s\"", a.State, gotApp.State)
	}
	if gotApp.PipelineEnriched != a.PipelineEnriched {
		t.Fatalf("expected \"%s\" got \"%s\"", a.PipelineEnriched, gotApp.PipelineEnriched)
	}
	if gotApp.PipelineFilenames != a.PipelineFilenames {
		t.Fatalf("expected \"%s\" got \"%s\"", a.PipelineFilenames, gotApp.PipelineFilenames)
	}
	if gotApp.PipelineOriginal != a.PipelineOriginal {
		t.Fatalf("expected \"%s\" got \"%s\"", a.PipelineOriginal, gotApp.PipelineOriginal)
	}
	if gotApp.ApplicationSpec != a.ApplicationSpec {
		t.Fatalf("expected \"%s\" got \"%s\"", a.ApplicationSpec, gotApp.ApplicationSpec)
	}
	if gotApp.Created.String() != a.Created.String() {
		t.Fatalf("expected \"%s\" got \"%s\"", a.Created.String(), gotApp.Created.String())
	}
	if gotApp.Updated.String() != a.Updated.String() {
		t.Fatalf("expected \"%s\" got \"%s\"", a.Updated.String(), gotApp.Updated.String())
	}
}
