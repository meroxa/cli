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

package global

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/meroxa/meroxa-go"
	"golang.org/x/oauth2"
)

const (
	clientID         = "2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb"
	meroxaBaseAPIURL = "https://api.meroxa.io"
)

func GetMeroxaAPIURL() string {
	if v := Config.GetString("API_URL"); v != "" {
		return v
	}

	return meroxaBaseAPIURL
}

func GetCLIUserInfo() (actor, actorUUID string, err error) {
	// Require login
	_, _, err = GetUserToken()

	/*
		 	We don't report client issues to the customer as it'll likely require `meroxa login` for any command.
			There are command that don't require client such as `meroxa env`, and we wouldn't like to throw an error,
			just because we can't emit events.
	*/
	if err != nil {
		return "", "", nil
	}

	// fetch actor account.
	actor = Config.GetString("ACTOR")
	actorUUID = Config.GetString("ACTOR_UUID")

	if actor == "" || actorUUID == "" {
		// call api to fetch
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // nolint:gomnd
		defer cancel()

		m, err := NewClient()

		if err != nil {
			return "", "", fmt.Errorf("meroxa: could not create Meroxa client: %v", err)
		}

		account, err := m.GetUser(ctx)

		if err != nil {
			return "", "", fmt.Errorf("meroxa: could not fetch Meroxa user: %v", err)
		}

		actor = account.Email
		actorUUID = account.UUID

		Config.Set("ACTOR", actor)
		Config.Set("ACTOR_UUID", actorUUID)

		err = Config.WriteConfig()

		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				err = Config.SafeWriteConfig()
			}
			if err != nil {
				return "", "", fmt.Errorf("meroxa: could not write config file: %v", err)
			}
		}
	}

	return actor, actorUUID, nil
}

func GetUserToken() (accessToken, refreshToken string, err error) {
	accessToken = Config.GetString("MEROXA_ACCESS_TOKEN")
	refreshToken = Config.GetString("MEROXA_REFRESH_TOKEN")
	if accessToken == "" && refreshToken == "" {
		// we need at least one token for creating an authenticated client
		return "", "", errors.New("please login or signup by running 'meroxa login'")
	}

	return accessToken, refreshToken, nil
}

func NewClient() (*meroxa.Client, error) {
	accessToken, refreshToken, err := GetUserToken()

	if err != nil {
		return nil, err
	}

	options := []meroxa.Option{
		meroxa.WithUserAgent(fmt.Sprintf("Meroxa CLI %s", Version)),
	}

	if flagDebug {
		options = append(options, meroxa.WithDumpTransport(os.Stdout))
	}
	if flagTimeout != 0 {
		options = append(options, meroxa.WithClientTimeout(flagTimeout))
	}
	if flagAPIURL != "" {
		options = append(options, meroxa.WithBaseURL(flagAPIURL))
	}

	// WithAuthentication needs to be added after WithDumpTransport
	// to catch requests to auth0
	options = append(options, meroxa.WithAuthentication(
		&oauth2.Config{
			ClientID: clientID,
			Endpoint: meroxa.OAuth2Endpoint,
		},
		accessToken,
		refreshToken,
		onTokenRefreshed,
	))

	return meroxa.New(options...)
}

// onTokenRefreshed tries to save the new token in the config.
func onTokenRefreshed(token *oauth2.Token) {
	Config.Set("MEROXA_ACCESS_TOKEN", token.AccessToken)
	Config.Set("MEROXA_REFRESH_TOKEN", token.RefreshToken)
	_ = Config.WriteConfig() // ignore error, it's a best effort
}
