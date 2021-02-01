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
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy a component",
	Long: `deprovision a component of the Meroxa platform, including pipelines,
 resources, connectors, functions etc...`,
}

var destroyResourceCmd = &cobra.Command{
	Use:   "resource <name>",
	Short: "destroy resource",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Resource Name
		resName := args[0]

		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// get Resource ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		res, err := c.GetResourceByName(ctx, resName)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		c, err = client()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		ctx = context.Background()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = c.DeleteResource(ctx, res.ID)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		prettyPrint("resource destroyed", res)
	},
}

var destroyConnectorCmd = &cobra.Command{
	Use:   "connector <name>",
	Short: "destroy connector",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Connector Name
		conName := args[0]

		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// get Connector ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		con, err := c.GetConnectorByName(ctx, conName)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		c, err = client()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		ctx = context.Background()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = c.DeleteConnector(ctx, con.ID)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		prettyPrint("connector destroyed", con)
	},
}

var destroyFunctionCmd = &cobra.Command{
	Use:   "function",
	Short: "destroy function",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy function called - Not Implemented")
	},
}

var destroyPipelineCmd = &cobra.Command{
	Use:   "pipeline <name>",
	Short: "destroy pipeline",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Pipeline Name
		pipelineName := args[0]

		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// get Pipeline ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pipeline, err := c.GetPipelineByName(ctx, pipelineName)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		c, err = client()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		ctx = context.Background()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = c.DeletePipeline(ctx, pipeline.ID)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		prettyPrint("pipeline destroyed", pipeline)
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	// Subcommands
	destroyCmd.AddCommand(destroyResourceCmd)
	destroyCmd.AddCommand(destroyConnectorCmd)
	destroyCmd.AddCommand(destroyFunctionCmd)
	destroyCmd.AddCommand(destroyPipelineCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
