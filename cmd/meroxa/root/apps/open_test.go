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
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"strings"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"github.com/meroxa/cli/cmd/meroxa/builder"
// 	basicMock "github.com/meroxa/cli/cmd/meroxa/global/mock"
// 	"github.com/meroxa/cli/log"
// 	"github.com/meroxa/cli/utils"
// )

// func TestOpenAppArgs(t *testing.T) {
// 	tests := []struct {
// 		args    []string
// 		err     error
// 		appName string
// 	}{
// 		{args: []string{"my-app-name"}, err: nil, appName: "my-app-name"},
// 	}

// 	for _, tt := range tests {
// 		cc := &Open{}
// 		err := cc.ParseArgs(tt.args)

// 		if err != nil && tt.err.Error() != err.Error() {
// 			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
// 		}

// 		if tt.appName != cc.args.NameOrUUID {
// 			t.Fatalf("expected \"%s\" got \"%s\"", tt.appName, cc.args.NameOrUUID)
// 		}
// 	}
// }

// func TestOpenAppFlags(t *testing.T) {
// 	expectedFlags := []struct {
// 		name      string
// 		required  bool
// 		shorthand string
// 		hidden    bool
// 	}{
// 		{name: "path", required: false},
// 	}

// 	c := builder.BuildCobraCommand(&Open{})

// 	for _, f := range expectedFlags {
// 		cf := c.Flags().Lookup(f.name)
// 		require.NotNil(t, cf)
// 		assert.Equal(t, f.shorthand, cf.Shorthand)
// 		assert.Equal(t, f.required, utils.IsFlagRequired(cf))
// 		assert.Equal(t, f.hidden, cf.Hidden)
// 	}
// }

// type mockOpener struct {
// 	startURL string
// }

// func (m *mockOpener) Start(URL string) error {
// 	m.startURL = URL
// 	return nil
// }

// func TestOpenAppExecution(t *testing.T) {
// 	ctx := context.Background()

// 	testCases := []struct {
// 		desc      string
// 		appArg    string
// 		tenant    string
// 		appFlag   string
// 		appPath   string
// 		expectURL string
// 		apiURL    string
// 		wantErr   error
// 	}{
// 		{
// 			desc:      "Successfully open app link with arg",
// 			appArg:    "app-name",
// 			tenant:    "test",
// 			expectURL: "https://test.na1.meroxa.cloud/apps/b0p2ok0dcjisn4z/detail",
// 			appPath:   "",
// 			apiURL:    "https://test.na1.meroxa.cloud",
// 		},
// 		{
// 			desc:    "Fail with bad path",
// 			appFlag: os.TempDir(),
// 			wantErr: errors.New("supply either ID/Name argument or --path flag"),
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			client := basicMock.NewMockBasicClient(ctrl)
// 			logger := log.NewTestLogger()
// 			os.Setenv("UNIT_TEST", "true")
// 			os.Setenv("MEROXA_API_URL", tc.apiURL)

// 			filter := &url.Values{}
// 			filter.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", tc.appArg, tc.appArg))

// 			httpResp := &http.Response{
// 				Body:       io.NopCloser(strings.NewReader(body)),
// 				Status:     "200 OK",
// 				StatusCode: 200,
// 			}
// 			if tc.wantErr == nil {
// 				client.EXPECT().CollectionRequest(ctx, "GET", collectionName, "", nil, *filter).Return(
// 					httpResp,
// 					nil,
// 				)
// 			}

// 			opener := &mockOpener{}
// 			o := &Open{
// 				logger: logger,
// 				Opener: opener,
// 				args: struct {
// 					NameOrUUID string
// 				}{NameOrUUID: tc.appArg},
// 				path:   tc.appPath,
// 				client: client,
// 			}

// 			err := o.Execute(ctx)

// 			if tc.wantErr != nil {
// 				require.Error(t, err)
// 				require.Contains(t, err.Error(), tc.wantErr.Error())
// 			}
// 			require.Equal(t, opener.startURL, tc.expectURL)
// 		})
// 	}
// }
