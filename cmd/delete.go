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
	"time"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a component",
	Long: `Deprovision a component of the Meroxa platform, including pipelines,
 resources, connectors, functions, etc...`,
}

var deleteResourceCmd = &cobra.Command{
	Use:   "resource <name>",
	Short: "Delete resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires resource name\n\nUsage:\n  meroxa delete resource <name>")
		}
		// Resource Name
		resName := args[0]

		c, err := client()
		if err != nil {
			return err
		}

		// get Resource ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = c.DeleteResource(ctx, res.ID)
		if err != nil {
			return err
		}

		prettyPrint("resource deleted", res)
		return nil
	},
}

var deleteConnectorCmd = &cobra.Command{
	Use:   "connector <name>",
	Short: "Delete connector",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires connector name\n\nUsage:\n  meroxa delete connector <name>")
		}

		// Connector Name
		conName := args[0]

		c, err := client()
		if err != nil {
			return err
		}

		// get Connector ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = c.DeleteConnector(ctx, con.ID)
		if err != nil {
			return err
		}

		prettyPrint("connector deleted", con)
		return nil
	},
}

var deletePipelineCmd = &cobra.Command{
	Use:   "pipeline <name>",
	Short: "Delete pipeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires pipeline name\n\nUsage:\n  meroxa delete pipeline <name>")
		}

		// Pipeline Name
		pipelineName := args[0]

		c, err := client()
		if err != nil {
			return err
		}

		// get Pipeline ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = c.DeletePipeline(ctx, pipeline.ID)
		if err != nil {
			return err
		}

		prettyPrint("Pipeline deleted", pipeline)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	// Subcommands
	deleteCmd.AddCommand(deleteResourceCmd)
	deleteCmd.AddCommand(deleteConnectorCmd)
	deleteCmd.AddCommand(deletePipelineCmd)
}
