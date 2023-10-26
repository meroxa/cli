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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/meroxa/cli/log"

	"github.com/golang/mock/gomock"
	basicMock "github.com/meroxa/cli/cmd/meroxa/global/mock"
)

func TestListApplicationExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := basicMock.NewMockBasicClient(ctrl)
	logger := log.NewTestLogger()

	appTime := AppTime{}
	appTime.UnmarshalJSON([]byte(`"2023-10-25 22:40:21.297Z"`))
	a := &Application{}
	a.Name = "my-env"
	a.State = "running"
	a.SpecVersion = "0.2.0"
	a.Created = appTime
	a.Updated = appTime

	a2 := &Application{}
	a2.Name = "my-env2"
	a2.State = "creating"
	a2.SpecVersion = "0.2.0"
	a2.Created = appTime
	a2.Updated = appTime

	allApps := []Application{*a, *a2}

	body := `{
		"page":1,
		"perPage":30,
		"totalItems":2,
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
		   },
		   {
			"collectionId":"gnhz55oi6tulkvs",
			"collectionName":"apps",
			"created":"2023-10-25 22:40:21.297Z",
			"id":"b0p2ok0dcjisn4z",
			"name":"my-env2",
			"specVersion":"0.2.0",
			"state":"creating",
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
	filter := &url.Values{}
	filter.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", a.Name, a.Name))
	output := &Applications{}

	httpResp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(body)),
		Status:     "200 OK",
		StatusCode: 200,
	}
	client.EXPECT().CollectionRequest(ctx, "GET", collectionName, "", nil, nil, output).Return(
		httpResp,
		nil,
	)

	list := &List{
		client: client,
		logger: logger,
	}

	err := list.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotJSONOutput := logger.JSONOutput()

	var gotApp Applications
	err = json.Unmarshal([]byte(gotJSONOutput), &gotApp)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	for i, app := range gotApp.Items {
		if app.Name != allApps[i].Name {
			t.Fatalf("expected \"%s\" got \"%s\"", a.Name, app.Name)
		}
		if app.SpecVersion != allApps[i].SpecVersion {
			t.Fatalf("expected \"%s\" got \"%s\"", a.SpecVersion, app.SpecVersion)
		}
		if app.State != allApps[i].State {
			t.Fatalf("expected \"%s\" got \"%s\"", a.State, app.State)
		}
		if app.Created.String() != allApps[i].Created.String() {
			t.Fatalf("expected \"%s\" got \"%s\"", a.Created.String(), app.Created.String())
		}
		if app.Updated.String() != allApps[i].Updated.String() {
			t.Fatalf("expected \"%s\" got \"%s\"", a.Updated.String(), app.Updated.String())
		}

	}
}
