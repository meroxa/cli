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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
)

func TestOpenAppArgs(t *testing.T) {
	tests := []struct {
		args    []string
		err     error
		appName string
	}{
		{args: []string{"my-app-name"}, err: nil, appName: "my-app-name"},
	}

	for _, tt := range tests {
		cc := &Open{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.appName != cc.args.NameOrUUID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.appName, cc.args.NameOrUUID)
		}
	}
}

func TestOpenAppFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "path", required: false},
	}

	c := builder.BuildCobraCommand(&Open{})

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

func TestOpenAppExecution(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		desc         string
		appArg       string
		appFlag      string
		errSubstring string
	}{
		{
			desc:   "Successfully open app link with arg",
			appArg: "app-name",
		},
		{
			desc: "Successfully open app link with flag",
		},
		{
			desc:         "Fail with bad path",
			appFlag:      "/tmp",
			errSubstring: "could not find an app.json file on path",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			os.Setenv("UNIT_TEST", "1")
			path := filepath.Join("/tmp", uuid.New().String())
			logger := log.NewTestLogger()
			cc := &Init{}
			cc.Logger(logger)
			cc.flags.Path = path
			cc.flags.Lang = "golang"
			if tc.appArg != "" {
				cc.args.appName = tc.appArg
				if tc.appFlag == "" {
					tc.appFlag = path + "/" + tc.appArg
				}
			} else {
				cc.args.appName = "my-app"
				if tc.appFlag == "" {
					tc.appFlag = path + "/my-app"
				}
			}

			err := cc.Execute(context.Background())
			if err != nil {
				t.Fatalf("unexpected error \"%s\"", err)
			}

			o := &Open{
				logger: logger,
			}
			if tc.appArg != "" {
				o.args = struct {
					NameOrUUID string
				}{NameOrUUID: tc.appArg}
			} else if tc.appFlag != "" {
				o.flags = struct {
					Path string `long:"path" usage:"Path to the app directory (default is local directory)"`
				}{Path: tc.appFlag}
			}

			err = o.Execute(ctx)
			if tc.errSubstring == "" && err != nil {
				t.Fatalf("not expected error, got \"%s\"", err.Error())
			} else if err != nil && !strings.Contains(err.Error(), tc.errSubstring) {
				t.Fatalf("failed to find expected error output(%s):\n%s", tc.errSubstring, err.Error())
			}
		})
	}
}
