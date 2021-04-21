/*
Copyright © 2020 Meroxa Inc

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

package global

import (
	"fmt"
	"os"

	"github.com/meroxa/meroxa-go"
	"github.com/meroxa/meroxa-go/auth"
	"golang.org/x/oauth2"
)

const (
	clientID = "2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb"
)

func NewClient() (*meroxa.Client, error) {
	options := []meroxa.Option{
		meroxa.WithUserAgent(fmt.Sprintf("Meroxa CLI %s", Version)),
	}

	if flagDebug {
		options = append(options, meroxa.WithDumpTransport(os.Stdout))
	}
	if flagTimeout != 0 {
		options = append(options, meroxa.WithClientTimeout(flagTimeout))
	}
	if flagApiUrl != "" {
		options = append(options, meroxa.WithBaseURL(flagApiUrl))
	}

	// WithAuthentication needs to be added after WithDumpTransport
	// to catch requests to auth0
	options = append(options, meroxa.WithAuthentication(
		&oauth2.Config{
			ClientID: clientID,
			Endpoint: auth.Endpoint,
		},
		Config.GetString("ACCESS_TOKEN"),
		Config.GetString("REFRESH_TOKEN"),
		onTokenRefreshed,
	))

	return meroxa.New(options...)
}

// onTokenRefreshed tries to save the new token in the config.
func onTokenRefreshed(token *oauth2.Token) {
	Config.Set("ACCESS_TOKEN", token.AccessToken)
	Config.Set("REFRESH_TOKEN", token.RefreshToken)
	_ = Config.WriteConfig() // ignore error, it's a best effort
}
