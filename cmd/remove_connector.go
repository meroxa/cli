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
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

// RemoveConnectorCmd represents the `meroxa remove connector` command
type RemoveConnector struct {
	name      string
	removeCmd *Remove
}

// RemoveConnectorClient represents the interface for meroxa client
type RemoveConnectorClient interface {
	GetConnectorByName(ctx context.Context, name string) (*meroxa.Connector, error)
	DeleteConnector(ctx context.Context, id int) error
}

func (rc *RemoveConnector) setArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires connector name\n\nUsage:\n  meroxa remove connector NAME")
	}
	// Connector Name
	rc.name = args[0]

	rc.removeCmd.componentType = "connector"
	rc.removeCmd.confirmableName = rc.name

	return nil
}

func (rc *RemoveConnector) execute(ctx context.Context, c RemoveConnectorClient) (*meroxa.Connector, error) {
	con, err := c.GetConnectorByName(ctx, rc.name)
	if err != nil {
		return nil, err
	}

	return con, c.DeleteConnector(ctx, con.ID)
}

func (rc *RemoveConnector) output(c *meroxa.Connector) {
	if flagRootOutputJSON {
		utils.JSONPrint(c)
	} else {
		fmt.Printf("connector %s successfully removed\n", c.Name)
	}
}

// command represents the `meroxa remove connector` command
func (rc *RemoveConnector) command() *cobra.Command {
	return &cobra.Command{
		Use:   "connector NAME",
		Short: "Remove connector",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return rc.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()
			if err != nil {
				return err
			}

			con, err := rc.execute(ctx, c)

			if err != nil {
				return err
			}

			rc.output(con)
			return nil
		},
	}
}
