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

	"github.com/meroxa/cli/cmd/meroxa/builder"

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

func TestListResourcesFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{name: "types", required: false, shorthand: ""},
	}

	c := builder.BuildCobraCommand(&List{})

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		}

		if f.shorthand != cf.Shorthand {
			t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
		}

		if f.required && !utils.IsFlagRequired(cf) {
			t.Fatalf("expected flag \"%s\" to be required", f.name)
		}
	}
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

func TestListResourceTypesExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockListResourcesClient(ctrl)
	logger := log.NewTestLogger()

	var types = []string{
		"postgres",
		"s3",
		"redshift",
		"mysql",
		"url",
		"mongodb",
		"elasticsearch",
		"snowflakedb",
		"bigquery",
	}

	client.
		EXPECT().
		ListResourceTypes(ctx).
		Return(types, nil)

	l := &List{
		client: client,
		logger: logger,
	}

	l.flags.Types = true

	err := l.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := utils.ResourceTypesTable(types)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotTypes []string
	err = json.Unmarshal([]byte(gotJSONOutput), &gotTypes)

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotTypes, types) {
		t.Fatalf("expected \"%v\", got \"%v\"", types, gotTypes)
	}
}
