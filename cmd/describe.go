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
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "describe a component",
	Long: `describe a component of the Meroxa data platform, including pipelines,
resources, connections, functions etc...`,
}

var describeResourceCmd = &cobra.Command{
	Use:   "resource <name>",
	Short: "describe resource",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := client()
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		res, err := c.GetResourceByName(ctx, name)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		prettyPrint("resource", res)
	},
}

var describeConnectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "describe connection",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client()
		intID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: ", err)
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		res, err := c.GetConnection(ctx, intID)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		prettyPrint("connection", res)
	},
}

var describeFunctionCmd = &cobra.Command{
	Use:   "function",
	Short: "describe function",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("describe function called - Not Implemented")
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)

	// Subcommands
	describeCmd.AddCommand(describeResourceCmd)
	describeCmd.AddCommand(describeConnectionCmd)
	describeCmd.AddCommand(describeFunctionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// describeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// describeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
