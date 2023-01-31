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

package resources

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/meroxa/cli/utils/display"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func getResources() []*meroxa.Resource {
	var resources []*meroxa.Resource
	r := utils.GenerateResource()
	return append(resources, &r)
}

func getResourcesWithEnvironment() []*meroxa.Resource {
	var resources []*meroxa.Resource
	r := utils.GenerateResourceWithEnvironment()
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
		} else {
			if f.shorthand != cf.Shorthand {
				t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
			}

			if f.required && !utils.IsFlagRequired(cf) {
				t.Fatalf("expected flag \"%s\" to be required", f.name)
			}
		}
	}
}

func TestListResourcesExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	resources := append(getResources(), getResourcesWithEnvironment()...)

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
	wantLeveledOutput := display.ResourcesTable(resources, false)

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
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	var types = []meroxa.ResourceType{
		{
			Name:         string(meroxa.ResourceTypePostgres),
			ReleaseStage: meroxa.ResourceTypeReleaseStageBeta,
			FormConfig: map[string]interface{}{
				meroxa.ResourceTypeFormConfigHumanReadableKey: "PostgreSQL",
			},
		},
		{
			Name:         string(meroxa.ResourceTypeS3),
			ReleaseStage: meroxa.ResourceTypeReleaseStageBeta,
			FormConfig: map[string]interface{}{
				meroxa.ResourceTypeFormConfigHumanReadableKey: "AWS S3",
			},
		},
		{
			Name:         string(meroxa.ResourceTypeRedshift),
			ReleaseStage: meroxa.ResourceTypeReleaseStageBeta,
			FormConfig: map[string]interface{}{
				meroxa.ResourceTypeFormConfigHumanReadableKey: "AWS Redshift",
			},
		},
		{
			Name:         string(meroxa.ResourceTypeMysql),
			ReleaseStage: meroxa.ResourceTypeReleaseStageBeta,
			FormConfig: map[string]interface{}{
				meroxa.ResourceTypeFormConfigHumanReadableKey: "MySQL",
			},
		},
		{
			Name:         string(meroxa.ResourceTypeUrl),
			ReleaseStage: meroxa.ResourceTypeReleaseStageBeta,
			FormConfig: map[string]interface{}{
				meroxa.ResourceTypeFormConfigHumanReadableKey: "Generic HTTP",
			},
		},
		{
			Name:         string(meroxa.ResourceTypeMongodb),
			ReleaseStage: meroxa.ResourceTypeReleaseStageBeta,
			FormConfig: map[string]interface{}{
				meroxa.ResourceTypeFormConfigHumanReadableKey: "MongoDB",
			},
		},
		{
			Name:         string(meroxa.ResourceTypeElasticsearch),
			ReleaseStage: meroxa.ResourceTypeReleaseStageBeta,
			FormConfig: map[string]interface{}{
				meroxa.ResourceTypeFormConfigHumanReadableKey: "Elasticsearch",
			},
		},
	}

	client.
		EXPECT().
		ListResourceTypesV2(ctx).
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
	wantLeveledOutput := display.ResourceTypesTable(types, false)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotTypes []meroxa.ResourceType
	err = json.Unmarshal([]byte(gotJSONOutput), &gotTypes)

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotTypes, types) {
		t.Fatalf("expected \"%v\", got \"%v\"", types, gotTypes)
	}
}
