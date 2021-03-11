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

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a component",
	Long: `Deprovision a component of the Meroxa platform, including pipelines,
 resources, and connectors`,
	SuggestFor: []string{"destroy", "delete"},
	Aliases:    []string{"rm", "delete"},
}

var removeEndpointCmd = &cobra.Command{
	Use:     "endpoint <name>",
	Aliases: []string{"endpoints"},
	Short:   "Remove endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires endpoint name\n\nUsage:\n  meroxa remove endpoint <name> [flags]")
		}
		name := args[0]

		c, err := client()
		if err != nil {
			return err
		}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
		defer cancel()

		return c.DeleteEndpoint(ctx, name)
	},
}

var removeResourceCmd = &cobra.Command{
	Use:   "resource <name>",
	Short: "Remove resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires resource name\n\nUsage:\n  meroxa remove resource <name>")
		}
		// Resource Name
		resName := args[0]

		c, err := client()
		if err != nil {
			return err
		}

		// get Resource ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
		defer cancel()

		res, err := c.GetResourceByName(ctx, resName)
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

		// TODO: Update meroxa-go to `RemoveResource` to match its implementation
		err = c.DeleteResource(ctx, res.ID)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(res)
		} else {
			fmt.Printf("Resource %s removed\n", res.Name)
		}
		return nil
	},
}

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

var removePipelineCmd = &cobra.Command{
	Use:   "pipeline <name>",
	Short: "Remove pipeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires pipeline name\n\nUsage:\n  meroxa remove pipeline <name>")
		}

		// Pipeline Name
		pipelineName := args[0]

		c, err := client()
		if err != nil {
			return err
		}

		// get Pipeline ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
		defer cancel()

		pipeline, err := c.GetPipelineByName(ctx, pipelineName)
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

		err = c.DeletePipeline(ctx, pipeline.ID)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(pipeline)
		} else {
			fmt.Printf("Pipeline %s removed\n", pipeline.Name)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)

	// Subcommands
	removeCmd.AddCommand(removeResourceCmd)
	removeCmd.AddCommand(removeConnectorCmd)
	removeCmd.AddCommand(removePipelineCmd)
	removeCmd.AddCommand(removeEndpointCmd)
}
