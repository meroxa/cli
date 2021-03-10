package cmd

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

import (
	"context"

	"github.com/meroxa/cli/display"
	"github.com/meroxa/meroxa-go"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List components",
	Long: `List the components of the Meroxa platform, including pipelines,
 resources, connectors, etc... You may also filter by type.`,
}

var listEndpointsCmd = &cobra.Command{
	Use:     "endpoint",
	Aliases: []string{"endpoints"},
	Short:   "List endpoints",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), meroxa.ClientTimeOut)
		defer cancel()

		ends, err := c.ListEndpoints(ctx)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(ends)
		} else {
			display.PrintEndpointsTable(ends)
		}

		return nil
	},
}

var listResourcesCmd = &cobra.Command{
	Use:     "resources",
	Short:   "List resources",
	Aliases: []string{"resource"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, meroxa.ClientTimeOut)
		defer cancel()

		rr, err := c.ListResources(ctx)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(rr)
		} else {
			display.PrintResourcesTable(rr)
		}
		return nil
	},
}

var listConnectorsCmd = &cobra.Command{
	Use:     "connectors",
	Short:   "List connectors",
	Aliases: []string{"connector"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, meroxa.ClientTimeOut)
		defer cancel()

		connectors, err := c.ListConnectors(ctx)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(connectors)
		} else {
			display.PrintConnectorsTable(connectors)
		}
		return nil
	},
}

var listResourceTypesCmd = &cobra.Command{
	Use:     "resource-types",
	Short:   "List resources-types",
	Aliases: []string{"resource-type"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, meroxa.ClientTimeOut)
		defer cancel()

		resTypes, err := c.ListResourceTypes(ctx)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(resTypes)
		} else {
			display.PrintResourceTypesTable(resTypes)
		}
		return nil
	},
}

var listPipelinesCmd = &cobra.Command{
	Use:     "pipelines",
	Short:   "List pipelines",
	Aliases: []string{"pipeline"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, meroxa.ClientTimeOut)
		defer cancel()

		rr, err := c.ListPipelines(ctx)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(rr)
		} else {
			display.PrintPipelinesTable(rr)
		}
		return nil
	},
}

var listTransformsCmd = &cobra.Command{
	Use:     "transforms",
	Short:   "List transforms",
	Aliases: []string{"transform"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, meroxa.ClientTimeOut)
		defer cancel()

		rr, err := c.ListTransforms(ctx)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(rr)
		} else {
			display.PrintTransformsTable(rr)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Subcommands
	listCmd.AddCommand(listResourcesCmd)
	listCmd.AddCommand(listConnectorsCmd)
	listCmd.AddCommand(listResourceTypesCmd)
	listCmd.AddCommand(listPipelinesCmd)
	listCmd.AddCommand(listTransformsCmd)
	listCmd.AddCommand(listEndpointsCmd)
}
