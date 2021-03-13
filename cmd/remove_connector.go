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
	"github.com/meroxa/cli/display"
	"github.com/spf13/cobra"
)

var removeConnectorCmd = &cobra.Command{
	Use:   "connector <name>",
	Short: "Remove connector",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires connector name\n\nUsage:\n  meroxa remove connector <name>")
		}

		// Connector Name
		conName := args[0]

		c, err := client()
		if err != nil {
			return err
		}

		// get Connector ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
		defer cancel()

		con, err := c.GetConnectorByName(ctx, conName)
		if err != nil {
			return err
		}

		c, err = client()
		if err != nil {
			return err
		}

		ctx = context.Background()
		ctx, cancel = context.WithTimeout(ctx, clientTimeOut)
		defer cancel()

		err = c.DeleteConnector(ctx, con.ID)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(con)
		} else {
			fmt.Printf("Connection %s removed\n", con.Name)
		}
		return nil
	},
}

func init() {
	removeCmd.AddCommand(removeConnectorCmd)
}
