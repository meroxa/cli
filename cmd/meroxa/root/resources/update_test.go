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
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/meroxa/cli/log"

	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/meroxa-go"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/utils"
)

func TestUpdateResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires resource name"), name: ""},
		{args: []string{"resource-name"}, err: nil, name: "resource-name"},
	}

	for _, tt := range tests {
		cc := &Update{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestUpdateResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "name", required: false, shorthand: ""},
		{name: "url", required: false, shorthand: "u"},
		{name: "metadata", required: false, shorthand: "m"},

		{name: "username", required: false, shorthand: ""},
		{name: "password", required: false, shorthand: ""},
		{name: "ca-cert", required: false, shorthand: ""},
		{name: "client-cert", required: false, shorthand: ""},
		{name: "client-key", required: false, shorthand: ""},
		{name: "ssl", required: false, shorthand: ""},
	}

	c := builder.BuildCobraCommand(&Update{})

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

		if cf.Hidden != f.hidden {
			if cf.Hidden {
				t.Fatalf("expected flag \"%s\" not to be hidden", f.name)
			} else {
				t.Fatalf("expected flag \"%s\" to be hidden", f.name)
			}
		}
	}
}

func TestUpdateResourceExecutionNoFlags(t *testing.T) {
	ctx := context.Background()
	u := &Update{}

	err := u.Execute(ctx)

	expected := "requires either `--name`, `--url`, `--metadata` or one of the credential flags"

	if err != nil && err.Error() != expected {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestUpdateResourceExecutionWithNewName(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateResourceClient(ctrl)
	logger := log.NewTestLogger()

	r := utils.GenerateResource()

	newName := "my-new-resource-name"
	nr := meroxa.UpdateResourceInput{
		Name: newName,
	}

	client.
		EXPECT().
		UpdateResource(ctx, r.Name, nr).
		Return(&r, nil)

	u := &Update{
		client: client,
		logger: logger,
	}

	u.args.Name = r.Name
	u.flags.Name = newName

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating resource %q...
Resource %q successfully updated!
`, u.args.Name, u.args.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotResource meroxa.Resource
	err = json.Unmarshal([]byte(gotJSONOutput), &gotResource)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotResource, r) {
		t.Fatalf("expected \"%v\", got \"%v\"", r, gotResource)
	}
}

func TestUpdateResourceExecutionWithNewMetadata(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateResourceClient(ctrl)
	logger := log.NewTestLogger()

	r := utils.GenerateResource()
	newMetadata := `{"metakey":"metavalue"}`

	var metadata map[string]interface{}

	_ = json.Unmarshal([]byte(newMetadata), &metadata)
	nr := meroxa.UpdateResourceInput{
		Metadata: metadata,
	}

	client.
		EXPECT().
		UpdateResource(ctx, r.Name, nr).
		Return(&r, nil)

	u := &Update{
		client: client,
		logger: logger,
	}

	u.args.Name = r.Name
	u.flags.Metadata = newMetadata

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating resource %q...
Resource %q successfully updated!
`, u.args.Name, u.args.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotResource meroxa.Resource
	err = json.Unmarshal([]byte(gotJSONOutput), &gotResource)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotResource, r) {
		t.Fatalf("expected \"%v\", got \"%v\"", r, gotResource)
	}
}

func TestUpdateResourceExecutionWithNewURL(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateResourceClient(ctrl)
	logger := log.NewTestLogger()

	r := utils.GenerateResource()
	newURL := "https://newUrl.io"

	nr := meroxa.UpdateResourceInput{
		URL: newURL,
	}

	client.
		EXPECT().
		UpdateResource(ctx, r.Name, nr).
		Return(&r, nil)

	u := &Update{
		client: client,
		logger: logger,
	}

	u.args.Name = r.Name
	u.flags.URL = newURL

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating resource %q...
Resource %q successfully updated!
`, u.args.Name, u.args.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotResource meroxa.Resource
	err = json.Unmarshal([]byte(gotJSONOutput), &gotResource)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotResource, r) {
		t.Fatalf("expected \"%v\", got \"%v\"", r, gotResource)
	}
}

func TestUpdateResourceExecutionWithNewCredentials(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateResourceClient(ctrl)
	logger := log.NewTestLogger()

	newUsername := "newUsername"

	r := utils.GenerateResource()

	// Updating one of their values only
	newCred := meroxa.Credentials{Username: newUsername}

	nr := meroxa.UpdateResourceInput{
		Credentials: &newCred,
	}

	client.
		EXPECT().
		UpdateResource(ctx, r.Name, nr).
		Return(&r, nil)

	u := &Update{
		client: client,
		logger: logger,
	}

	u.args.Name = r.Name
	u.flags.Username = newUsername

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating resource %q...
Resource %q successfully updated!
`, u.args.Name, u.args.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotResource meroxa.Resource
	err = json.Unmarshal([]byte(gotJSONOutput), &gotResource)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotResource, r) {
		t.Fatalf("expected \"%v\", got \"%v\"", r, gotResource)
	}
}
