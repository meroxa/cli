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
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		require.NotNil(t, cf)
		assert.Equal(t, f.shorthand, cf.Shorthand)
		assert.Equal(t, f.required, utils.IsFlagRequired(cf))
		assert.Equal(t, f.hidden, cf.Hidden)
	}
}

type mockOpener struct {
	startURL string
}

func (m *mockOpener) Start(URL string) error {
	m.startURL = URL
	return nil
}

func TestOpenAppExecution(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		desc      string
		appArg    string
		appFlag   string
		expectURL string
		wantErr   error
	}{
		{
			desc:      "Successfully open app link with arg",
			appArg:    "app-name",
			expectURL: "https://dashboard.meroxa.io/apps/app-name/detail",
		},
		{
			desc:      "Successfully open app link with flag",
			expectURL: "https://dashboard.meroxa.io/apps/my-app/detail",
		},
		{
			desc:    "Fail with bad path",
			appFlag: os.TempDir(),
			wantErr: errors.New("could not find an app.json file on path"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			os.Setenv("UNIT_TEST", "1")
			path := filepath.Join(os.TempDir(), uuid.New().String())
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
			require.NoError(t, err)

			opener := &mockOpener{}
			o := &Open{
				logger: logger,
				Opener: opener,
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
			if tc.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr.Error())
			}
			require.Equal(t, opener.startURL, tc.expectURL)
		})
	}
}
