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
	"github.com/meroxa/cli/display"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
	"time"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a component",
	Long:  `Describe a component of the Meroxa data platform, including resources and connectors`,
}

var describeResourceCmd = &cobra.Command{
	Use:   "resource <name>",
	Short: "Describe resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires resource name\n\nUsage:\n  meroxa describe resource <name> [flags]")
		}
		name := args[0]

		c, err := client()
		if err != nil {
			return err
		}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		res, err := c.GetResourceByName(ctx, name)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(res)
		} else {
			display.PrettyPrint("resource", res)
		}
		return nil
	},
}

var describeConnectorCmd = &cobra.Command{
	Use:   "connector [name]",
	Short: "Describe connector",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires connector name\n\nUsage:\n  meroxa describe connector <name> [flags]")
		}
		var (
			err  error
			conn *meroxa.Connector
		)
		name := args[0]
		c, err := client()
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		conn, err = c.GetConnectorByName(ctx, name)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(conn)
		} else {
			display.PrettyPrint("connector", conn)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(describeCmd)

	// Subcommands
	describeCmd.AddCommand(describeResourceCmd)
	describeCmd.AddCommand(describeConnectorCmd)
}
