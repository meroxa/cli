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
	"io"
	"net/http"
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
	err := appTime.UnmarshalJSON([]byte(`"2024-04-01 20:13:20.111Z"`))
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	a := &Application{
		Name:    "test-pipeline-1",
		State:   "provisioned",
		Created: appTime,
		Updated: appTime,
	}

	a2 := &Application{
		Name:    "test-pipeline-1",
		State:   "provisioned",
		Created: appTime,
		Updated: appTime,
	}

	allApps := []Application{*a, *a2}

	httpResp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(body)),
		Status:     "200 OK",
		StatusCode: 200,
	}
	client.EXPECT().CollectionRequest(ctx, "GET", applicationCollection, "", nil, nil).Return(
		httpResp,
		nil,
	)

	list := &List{
		client: client,
		logger: logger,
	}

	err = list.Execute(ctx)
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
