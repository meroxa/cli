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

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"testing"

// 	"github.com/google/uuid"
// 	basicMock "github.com/meroxa/cli/cmd/meroxa/global/mock"
// 	"github.com/stretchr/testify/require"

// 	"github.com/golang/mock/gomock"

// 	"github.com/meroxa/cli/log"
// )

// const (
// 	body = `
// 	{
// 		"page": 1,
// 		"perPage": 30,
// 		"totalItems": 1,
// 		"totalPages": 1,
// 		"items": [
// 		  {
// 			"collectionId": "77byam8idl1rv8b",
// 			"collectionName": "conduitapps",
// 			"config": null,
// 			"created": "2024-04-01 20:13:20.111Z",
// 			"deployment_id": [
// 			  "tcsmsunmfo5v8kw"
// 			],
// 			"id": "lxjcdlsvet3aeoe",
// 			"name": "test-pipeline-1",
// 			"pipeline_enriched": "version: \"2.2\"\npipelines:\n    - id: cp-pipeline-generator-log-source-generator.0\n      status: \"\"\n      name: pipeline-generator-log-name-0\n      description: \"\"\n      connectors:\n        - id: source-generator.0\n          type: source\n          plugin: builtin:generator\n          name: \"\"\n          settings:\n            format.options: event_id:int,pg_generator:bool,sensor_id:int,msg:string,triggered:bool\n            format.type: structured\n            readTime: 1s\n          processors: []\n        - id: kafka-source-generator.0\n          type: destination\n          plugin: builtin:kafka\n          name: kafka-source-generator.0.0\n          settings:\n            servers: 127.0.0.1:19092\n            topic: default.a9fc5274-c1df-4e85-a15a-a71337291817.0\n          processors: []\n      processors: []\n      dead-letter-queue:\n        plugin: builtin:log\n        settings:\n            level: warn\n            message: record delivery failed\n        window-size: 1\n        window-nack-threshold: 0\n    - id: cp-pipeline-generator-log-log-destination.1\n      status: \"\"\n      name: pipeline-generator-log-name-1\n      description: \"\"\n      connectors:\n        - id: log-destination.1\n          type: destination\n          plugin: builtin:log\n          name: \"\"\n          settings: {}\n          processors: []\n        - id: kafka-log-destination.1\n          type: source\n          plugin: builtin:kafka\n          name: kafka-log-destination.1.1\n          settings:\n            servers: 127.0.0.1:19092\n            topic: default.a9fc5274-c1df-4e85-a15a-a71337291817.0\n          processors: []\n      processors: []\n      dead-letter-queue:\n        plugin: builtin:log\n        settings:\n            level: warn\n            message: record delivery failed\n        window-size: 1\n        window-nack-threshold: 0\n",
// 			"pipeline_filename": "test-pipeline-1.yaml",
// 			"pipeline_original": "version: \"2.2\"\npipelines:\n    - id: pipeline-generator-log\n      status: running\n      name: pipeline-generator-log-name\n      description: \"\"\n      connectors:\n        - id: source-generator\n          type: source\n          plugin: builtin:generator\n          name: \"\"\n          settings:\n            format.options: event_id:int,pg_generator:bool,sensor_id:int,msg:string,triggered:bool\n            format.type: structured\n            readTime: 1s\n          processors: []\n        - id: log-destination\n          type: destination\n          plugin: builtin:log\n          name: \"\"\n          settings: {}\n          processors: []\n      processors: []\n      dead-letter-queue:\n        plugin: \"\"\n        settings: {}\n        window-size: null\n        window-nack-threshold: null\n",
// 			"state": "provisioned",
// 			"stream_tech": "kafka",
// 			"updated": "2024-04-01 20:13:20.111Z"
// 		  }
// 		]
// 	  }

// 	`
// )

// func TestDescribeApplicationArgs(t *testing.T) {
// 	tests := []struct {
// 		args []string
// 		err  error
// 		name string
// 	}{
// 		{args: nil, err: errors.New("requires app name or UUID"), name: ""},
// 		{args: []string{"ApplicationName"}, err: nil, name: "ApplicationName"},
// 	}

// 	for _, tt := range tests {
// 		ar := &Describe{}
// 		err := ar.ParseArgs(tt.args)

// 		if err != nil && tt.err.Error() != err.Error() {
// 			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
// 		}

// 		if tt.name != ar.args.nameOrUUID {
// 			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ar.args.nameOrUUID)
// 		}
// 	}
// }

// func TestDescribeApplicationExecution(t *testing.T) {
// 	ctx := context.Background()
// 	ctrl := gomock.NewController(t)
// 	client := basicMock.NewMockBasicClient(ctrl)
// 	logger := log.NewTestLogger()

// 	path := filepath.Join(os.TempDir(), uuid.NewString())
// 	appName := "test-pipeline-1"
// 	appTime := AppTime{}
// 	err := appTime.UnmarshalJSON([]byte(`"2024-04-01 20:13:20.111Z"`))
// 	if err != nil {
// 		t.Fatalf("not expected error, got \"%s\"", err.Error())
// 	}

// 	i := &Init{
// 		logger: logger,
// 		args:   struct{ appName string }{appName: appName},
// 		flags: struct {
// 			Lang        string "long:\"lang\" short:\"l\" usage:\"language to use (js|go|py)\" required:\"true\""
// 			Path        string "long:\"path\" usage:\"path where application will be initialized (current directory as default)\""
// 			ModVendor   bool   "long:\"mod-vendor\" usage:\"whether to download modules to vendor or globally while initializing a Go application\""
// 			SkipModInit bool   "long:\"skip-mod-init\" usage:\"whether to run 'go mod init' while initializing a Go application\""
// 		}{
// 			Path:        path,
// 			ModVendor:   false,
// 			SkipModInit: true,
// 		},
// 	}

// 	a := &Application{
// 		Name:    appName,
// 		State:   "running",
// 		Created: appTime,
// 		Updated: appTime,
// 	}

// 	err = i.Execute(ctx)
// 	defer func(path string) {
// 		os.RemoveAll(path)
// 	}(path)
// 	require.NoError(t, err)

// 	filter := &url.Values{}
// 	filter.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", a.Name, a.Name))

// 	httpResp := &http.Response{
// 		Body:       io.NopCloser(strings.NewReader(body)),
// 		Status:     "200 OK",
// 		StatusCode: 200,
// 	}
// 	client.EXPECT().CollectionRequest(ctx, "GET", applicationCollection, "", nil, *filter).Return(
// 		httpResp,
// 		nil,
// 	)

// 	client.EXPECT().AddHeader("Meroxa-CLI-App-Version", gomock.Any()).Times(1)

// 	dc := &Describe{
// 		client: client,
// 		logger: logger,
// 		args:   struct{ nameOrUUID string }{nameOrUUID: a.Name},
// 		flags: struct {
// 			Path string "long:\"path\" usage:\"Path to the app directory (default is local directory)\""
// 		}{Path: filepath.Join(path, appName)},
// 	}

// 	err = dc.Execute(ctx)
// 	if err != nil {
// 		t.Fatalf("not expected error, got %q", err.Error())
// 	}

// 	gotJSONOutput := logger.JSONOutput()

// 	var gotApp Application
// 	err = json.Unmarshal([]byte(gotJSONOutput), &gotApp)
// 	if err != nil {
// 		t.Fatalf("not expected error, got %q", err.Error())
// 	}

// 	if gotApp.Name != a.Name {
// 		t.Fatalf("expected \"%s\" got \"%s\"", a.Name, gotApp.Name)
// 	}

// 	if gotApp.State != a.State {
// 		t.Fatalf("expected \"%s\" got \"%s\"", a.State, gotApp.State)
// 	}
// 	if gotApp.Created.String() != a.Created.String() {
// 		t.Fatalf("expected \"%s\" got \"%s\"", a.Created.String(), gotApp.Created.String())
// 	}
// 	if gotApp.Updated.String() != a.Updated.String() {
// 		t.Fatalf("expected \"%s\" got \"%s\"", a.Updated.String(), gotApp.Updated.String())
// 	}
// }
