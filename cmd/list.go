/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"time"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list components",
	Long: `list the components of the Meroxa platform, including pipelines,
 resources, connections, functions etc... You may also filter by type.`,
}

var listResourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "list resources",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		rr, err := c.ListResources(ctx)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		switch output {
		case "json":
			prettyPrint("resources", rr)
		default:
			printResourcesTable(rr)
		}

	},
}

var listConnectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "list connections",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		connections, err := c.ListConnections(ctx)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		switch output {
		case "json":
			prettyPrint("connections", connections)
		default:
			printConnectionsTable(connections)
		}
	},
}

var listResourceTypesCmd = &cobra.Command{
	Use:   "resource-types",
	Short: "list resources-types",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		resTypes, err := c.ListResourceTypes(ctx)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		switch output {
		case "json":
			prettyPrint("resource types", resTypes)
		default:
			printResourceTypesTable(resTypes)
		}
	},
}

var listPipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "list pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		rr, err := c.ListPipelines(ctx)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		prettyPrint("pipelines", rr)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Subcommands
	listCmd.PersistentFlags().StringP("output", "o", "table", "output format [json|table]")
	listCmd.AddCommand(listResourcesCmd)
	listCmd.AddCommand(listConnectionsCmd)
	listCmd.AddCommand(listResourceTypesCmd)
	listCmd.AddCommand(listPipelinesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
