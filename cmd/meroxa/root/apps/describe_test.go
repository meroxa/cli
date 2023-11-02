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
	turbineMock "github.com/meroxa/cli/cmd/meroxa/turbine/mock"
	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"

	"github.com/meroxa/cli/log"
)

const (
	body = `{
		"page":1,
		"perPage":30,
		"totalItems":1,
		"totalPages":1,
		"items":[
		   {
			  "collectionId":"gnhz55oi6tulkvs",
			  "collectionName":"apps",
			  "created":"2023-10-25 22:40:21.297Z",
			  "id":"b0p2ok0dcjisn4z",
			  "name":"my-env",
			  "specVersion":"0.2.0",
			  "state":"running",
			  "updated":"2023-10-25 22:40:21.297Z",
			  "spec":{
				 "connectors":[
					{
					   "collection":"collection_name",
					   "resource":"source_name",
					   "type":"source",
					   "uuid":"5ce244be-e404-4fc1-b486-a35ee200fd27"
					},
					{
					   "collection":"collection_archive",
					   "resource":"destination_name",
					   "type":"destination",
					   "uuid":"0362c2df-6e99-445e-b95e-a798e69a651f"
					}
				 ],
				 "definition":{
					"git_sha":"f7baf1e05df0becdf946847b8f7411d22988a3d7\n",
					"metadata":{
					   "spec_version":"0.2.0",
					   "turbine":{
						  "language":"golang",
						  "version":"v2.1.3"
					   }
					}
				 },
				 "functions":[
					{
					   "image":"turbine-newgo.tar.gz",
					   "name":"anonymize",
					   "uuid":"04b0d690-dd44-4df3-8636-6f0c4dfb5c93"
					}
				 ],
				 "streams":[
					{
					   "from_uuid":"5ce244be-e404-4fc1-b486-a35ee200fd27",
					   "name":"5ce244be-e404-4fc1-b486-a35ee200fd27_04b0d690-dd44-4df3-8636-6f0c4dfb5c93",
					   "to_uuid":"04b0d690-dd44-4df3-8636-6f0c4dfb5c93",
					   "uuid":"ef1e3681-fbaa-4bff-9d21-6e010bcdec3e"
					},
					{
					   "from_uuid":"04b0d690-dd44-4df3-8636-6f0c4dfb5c93",
					   "name":"04b0d690-dd44-4df3-8636-6f0c4dfb5c93_0362c2df-6e99-445e-b95e-a798e69a651f",
					   "to_uuid":"0362c2df-6e99-445e-b95e-a798e69a651f",
					   "uuid":"06c89e49-753d-4a54-81f1-ee1e036003e6"
					}
				 ]
			  }
		   }
		]
	 }`
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
	mockTurbineCLI := turbineMock.NewMockCLI(ctrl)

	path := filepath.Join(os.TempDir(), uuid.NewString())
	appName := "my-env"
	appTime := AppTime{}
	err := appTime.UnmarshalJSON([]byte(`"2023-10-25 22:40:21.297Z"`))
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	i := &Init{
		logger: logger,
		args:   struct{ appName string }{appName: appName},
		flags: struct {
			Lang        string "long:\"lang\" short:\"l\" usage:\"language to use (js|go|py)\" required:\"true\""
			Path        string "long:\"path\" usage:\"path where application will be initialized (current directory as default)\""
			ModVendor   bool   "long:\"mod-vendor\" usage:\"whether to download modules to vendor or globally while initializing a Go application\""
			SkipModInit bool   "long:\"skip-mod-init\" usage:\"whether to run 'go mod init' while initializing a Go application\""
		}{
			Lang:        string(ir.GoLang),
			Path:        path,
			ModVendor:   false,
			SkipModInit: true,
		},
	}

	a := &Application{
		Name:        appName,
		State:       "running",
		SpecVersion: "0.2.0",
		Created:     appTime,
		Updated:     appTime,
	}

	err = i.Execute(ctx)
	defer func(path string) {
		os.RemoveAll(path)
	}(path)
	require.NoError(t, err)

	filter := &url.Values{}
	filter.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", a.Name, a.Name))

	httpResp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(body)),
		Status:     "200 OK",
		StatusCode: 200,
	}
	client.EXPECT().CollectionRequest(ctx, "GET", collectionName, "", nil, *filter).Return(
		httpResp,
		nil,
	)

	mockTurbineCLI.EXPECT().GetVersion(ctx).Return("1.0", nil)
	client.EXPECT().AddHeader("Meroxa-CLI-App-Lang", string(ir.GoLang)).Times(1)
	client.EXPECT().AddHeader("Meroxa-CLI-App-Version", gomock.Any()).Times(1)

	dc := &Describe{
		client:     client,
		logger:     logger,
		turbineCLI: mockTurbineCLI,
		args:       struct{ nameOrUUID string }{nameOrUUID: a.Name},
		flags: struct {
			Path string "long:\"path\" usage:\"Path to the app directory (default is local directory)\""
		}{Path: filepath.Join(path, appName)},
	}

	err = dc.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotJSONOutput := logger.JSONOutput()

	var gotApp Application
	err = json.Unmarshal([]byte(gotJSONOutput), &gotApp)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if gotApp.Name != a.Name {
		t.Fatalf("expected \"%s\" got \"%s\"", a.Name, gotApp.Name)
	}
	if gotApp.SpecVersion != a.SpecVersion {
		t.Fatalf("expected \"%s\" got \"%s\"", a.SpecVersion, gotApp.SpecVersion)
	}
	if gotApp.State != a.State {
		t.Fatalf("expected \"%s\" got \"%s\"", a.State, gotApp.State)
	}
	if gotApp.Created.String() != a.Created.String() {
		t.Fatalf("expected \"%s\" got \"%s\"", a.Created.String(), gotApp.Created.String())
	}
	if gotApp.Updated.String() != a.Updated.String() {
		t.Fatalf("expected \"%s\" got \"%s\"", a.Updated.String(), gotApp.Updated.String())
	}
}
