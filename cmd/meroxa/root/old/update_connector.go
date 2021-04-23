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

package old

import (
	"errors"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/utils"
	"github.com/spf13/cobra"
)

// UpdateConnectorCmd represents the `meroxa update connector` command.
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

			c, err := global.NewClient()
			if err != nil {
				return err
			}

			// call meroxa-go to update connector status with name
			if !FlagRootOutputJSON {
				fmt.Printf("Updating %s connector...\n", conName)
			}

			con, err := c.UpdateConnectorStatus(cmd.Context(), conName, state)
			if err != nil {
				return err
			}

			if FlagRootOutputJSON {
				utils.JSONPrint(con)
			} else {
				fmt.Printf("Connector %s successfully updated!\n", con.Name)
			}

			return nil
		},
	}

	// TODO: Validate state has to be either of pause|resume|restart
	updateConnectorCmd.Flags().StringVarP(&state, "state", "", "", "connector state")
	_ = updateConnectorCmd.MarkFlagRequired("state")

	return updateConnectorCmd
}
