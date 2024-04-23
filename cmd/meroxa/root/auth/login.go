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

package auth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"golang.org/x/term"
)

var (
	_ builder.CommandWithDocs        = (*Login)(nil)
	_ builder.CommandWithLogger      = (*Login)(nil)
	_ builder.CommandWithExecute     = (*Login)(nil)
	_ builder.CommandWithConfig      = (*Login)(nil)
	_ builder.CommandWithBasicClient = (*Login)(nil)
)

type Login struct {
	logger log.Logger
	config config.Config
	client global.BasicClient
}

func (l *Login) Usage() string {
	return "login"
}

func (l *Login) Docs() builder.Docs {
	return builder.Docs{
		Short: "Login to a Meroxa Platform tenant",
	}
}

func (l *Login) BasicClient(client global.BasicClient) {
	l.client = client
}

type authRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type pocketbaseResponse struct {
	Token  string                 `json:"token"`
	Record map[string]interface{} `json:"record"`
}

type pocketbaseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	Identity map[string]any `json:"identity"`
}

func (l *Login) getSetTenantSetting(globalSetting string, reader *bufio.Reader) (string, error) {
	prompted, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if globalSetting != "" {
		l.config.Set(globalSetting, strings.TrimSpace(prompted))
	}
	return strings.TrimSuffix(strings.TrimSpace(prompted), "\n"), nil
}

func (l *Login) settingsPromt(settingType, globalSetting, savedValue string, reader *bufio.Reader) error {
	if savedValue != "" {
		fmt.Printf("Found saved tenant %s, %s. Would you like to login using this tenant %s? (Y/N) \n", settingType, savedValue, settingType)
		yN, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		yN = strings.ToLower(yN)
		if strings.TrimSpace(yN) == "n" {
			fmt.Printf("Please enter value for %s: \n", settingType)
			_, err = l.getSetTenantSetting(globalSetting, reader)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Printf("Please enter value for %s: \n", settingType)
		_, err := l.getSetTenantSetting(globalSetting, reader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Login) Execute(ctx context.Context) error {
	reader := bufio.NewReader(os.Stdin)
	err := l.settingsPromt("url", global.TenantURL, global.GetMeroxaTenantURL(), reader)
	if err != nil {
		return err
	}
	l.logger.Info(ctx, fmt.Sprintf("Logging into tenant - %s", global.GetMeroxaTenantURL()))

	err = l.settingsPromt("email", global.TenantEmailAddress, global.GetMeroxaTenantUser(), reader)
	if err != nil {
		return err
	}

	fmt.Println("Please enter password: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		return err
	}

	req := authRequest{
		Identity: l.config.GetString(global.TenantEmailAddress),
		Password: strings.TrimSuffix(strings.TrimSpace(string(password)), "\n"),
	}

	var pbResp pocketbaseResponse
	var pbRespError pocketbaseError

	err = l.client.ResetBaseURL()
	if err != nil {
		return err
	}

	http, err := l.client.URLRequest(ctx, "POST", "/api/collections/users/auth-with-password", req, nil, nil)
	if err != nil {
		return err
	}
	if http.StatusCode == 404 {
		return fmt.Errorf("failed to log in to tenant - %s  . Please double check that this is a correct tenant url", global.GetMeroxaTenantURL())
	} else if http.StatusCode != 200 {
		err = json.NewDecoder(http.Body).Decode(&pbRespError)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to log in : %s %s %+v", pbRespError.Message, pbRespError.Data.Identity["message"], pbRespError)
	} else {
		err = json.NewDecoder(http.Body).Decode(&pbResp)
	}
	l.config.Set(global.AccessTokenEnv, pbResp.Token)
	l.config.Set(global.TenantEmailAddress, l.config.GetString(global.TenantEmailAddress))
	l.config.Set(global.TenantURL, l.config.GetString(global.TenantURL))

	l.logger.Info(ctx, fmt.Sprintf("\n ✨ Successfully logged into %s", global.GetMeroxaTenantURL()))
	return err
}

func (l *Login) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *Login) Config(cfg config.Config) {
	l.config = cfg
}
