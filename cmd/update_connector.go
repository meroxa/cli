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
	"errors"
	"fmt"
	"time"

	"github.com/meroxa/cli/utils"
	"github.com/spf13/cobra"
)

// UpdateConnectorCmd represents the `meroxa update connector` command
func UpdateConnectorCmd() *cobra.Command {
	var state string

	updateConnectorCmd := &cobra.Command{
		Use:     "connector NAME --state pause | resume | restart",
		Aliases: []string{"connectors"},
		Short:   "Update connector state",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires connector name\n\nUsage:\n  meroxa update connector NAME --state pause | resume | restart")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Connector Name
			conName := args[0]

			c, err := client()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			// call meroxa-go to update connector status with name
			if !flagRootOutputJSON {
				fmt.Printf("Updating %s connector...\n", conName)
			}

			con, err := c.UpdateConnectorStatus(ctx, conName, state)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				utils.JSONPrint(con)
			} else {
				fmt.Printf("Connector %s successfully updated!\n", con.Name)
			}

			return nil
		},
	}

	// TODO: Validate state has to be either of pause|resume|restart
	updateConnectorCmd.Flags().StringVarP(&state, "state", "", "", "connector state")
	updateConnectorCmd.MarkFlagRequired("state")

	return updateConnectorCmd
}
