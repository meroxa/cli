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

package global

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

// userInfoStale checks if user information was updated within a 24h period.
func userInfoStale() bool {
	updatedAt := Config.GetTime(UserInfoUpdatedAtEnv)
	if updatedAt.IsZero() {
		return true
	}

	duration := time.Now().UTC().Sub(updatedAt)
	return duration.Hours() > 24
}

func GetCLIUserInfo() (err error) {
	// Require login
	_, _, err = GetUserToken()

	/*
		 	We don't report client issues to the customer as it'll likely require `meroxa login` for any command.
			There are command that don't require client such as `meroxa config`, and we wouldn't like to throw an error,
			just because we can't emit events.
	*/
	if err != nil {
		return nil
	}

	if userInfoStale() {
		// call api to fetch
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		m, err := NewOauthClient()
		if err != nil {
			return fmt.Errorf("meroxa: could not create Meroxa client: %v", err)
		}

		user, err := m.GetUser(ctx)
		if err != nil {
			return fmt.Errorf("meroxa: could not fetch Meroxa user: %v", err)
		}

		// write existing feature flags enabled
		Config.Set(UserFeatureFlagsEnv, strings.Join(user.Features, " "))

		// write when was the last time we updated user info
		Config.Set(UserInfoUpdatedAtEnv, time.Now().UTC())

		err = Config.WriteConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				err = Config.SafeWriteConfig()
			}
			if err != nil {
				return fmt.Errorf("meroxa: could not write config file: %v", err)
			}
		}
	}

	return nil
}

func GetUserToken() (accessToken, refreshToken string, err error) {
	accessToken = Config.GetString(AccessTokenEnv)
	refreshToken = Config.GetString(RefreshTokenEnv)
	if accessToken == "" && refreshToken == "" {
		// we need at least one token for creating an authenticated client
		return "", "", errors.New("please login or signup by running 'meroxa login'")
	}

	return accessToken, refreshToken, nil
}

func NewOauthClient() (meroxa.Client, error) {
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
			ClientID: GetMeroxaAuthClientID(),
			Endpoint: oauthEndpoint(GetMeroxaAuthDomain()),
		},
		accessToken,
		refreshToken,
		onTokenRefreshed,
	))

	options = append(options, meroxa.WithHeader("Meroxa-CLI-Version", Version))
	return meroxa.New(options...)
}

func oauthEndpoint(domain string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("https://%s/authorize", domain),
		TokenURL: fmt.Sprintf("https://%s/oauth/token", domain),
	}
}

// onTokenRefreshed tries to save the new token in the config.
func onTokenRefreshed(token *oauth2.Token) {
	Config.Set(AccessTokenEnv, token.AccessToken)
	Config.Set(RefreshTokenEnv, token.RefreshToken)
	_ = Config.WriteConfig() // ignore error, it's a best effort
}
