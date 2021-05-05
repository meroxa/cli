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

package resources

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

func getResources() []*meroxa.Resource {
	var resources []*meroxa.Resource
	r := utils.GenerateResource()
	return append(resources, &r)
}

func TestListResourcesExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockListResourcesClient(ctrl)
	logger := log.NewTestLogger()

	resources := getResources()

	client.
		EXPECT().
		ListResources(ctx).
		Return(resources, nil)

	l := &List{
		client: client,
		logger: logger,
	}

	err := l.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := utils.ResourcesTable(resources)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotResources []meroxa.Resource
	err = json.Unmarshal([]byte(gotJSONOutput), &gotResources)

	var lr []meroxa.Resource

	for _, r := range resources {
		lr = append(lr, *r)
	}

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotResources, lr) {
		t.Fatalf("expected \"%v\", got \"%v\"", lr, gotResources)
	}
}
