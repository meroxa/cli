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

package cmd

import (
	"context"
	"fmt"

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
	var err error

	user, err := c.GetUser(ctx)
	return user, err
}

func (gu *GetUser) output(user *meroxa.User) {
	if flagRootOutputJSON {
		utils.JSONPrint(user)
	} else {
		fmt.Printf("%s\n", user.Email)
	}
}

// command represents the `meroxa whoami` command
func (gu *GetUser) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Display the current logged in user\n",
		Example: "\n" +
			"meroxa whoami'\n",

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()
			if err != nil {
				return err
			}

			u, err := gu.execute(ctx, c)

			if err != nil {
				return err
			}

			gu.output(u)

			return nil
		},
	}

	return cmd
}
