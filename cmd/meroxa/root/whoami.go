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

package root

import (
	"context"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/root/deprecated"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

type GetUserClient interface {
	GetUser(ctx context.Context) (*meroxa.User, error)
}

type GetUser struct {
}

func (gu *GetUser) execute(ctx context.Context, c GetUserClient) (*meroxa.User, error) {
	// TODO think about extracting the info from the access token.
	//  The access token is a JWT and contains these fields:
	//  * "https://api.meroxa.io/v1/email": "john.doe@example.com"
	//  * "https://api.meroxa.io/v1/username": "John Doe"
	var err error

	user, err := c.GetUser(ctx)
	return user, err
}

func (gu *GetUser) output(user *meroxa.User) {
	if deprecated.FlagRootOutputJSON {
		utils.JSONPrint(user)
	} else {
		fmt.Printf("%s\n", user.Email)
	}
}

// Command represents the `meroxa whoami` command.
func (gu *GetUser) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Display the current logged in user\n",
		Example: "\n" +
			"meroxa whoami'\n",

		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := global.NewClient()
			if err != nil {
				return err
			}

			u, err := gu.execute(cmd.Context(), c)

			if err != nil {
				return err
			}

			gu.output(u)

			return nil
		},
	}

	return cmd
}
