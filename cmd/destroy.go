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
	"fmt"

	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy a component",
	Long: `deprovision a component of the Meroxa platform, including pipelines,
 resources, connections, functions etc...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy called")
	},
}

var destroyResourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "destroy resource",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy resource called")
	},
}

var destroyConnectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "destroy connection",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy connection called")
	},
}

var destroyFunctionCmd = &cobra.Command{
	Use:   "function",
	Short: "destroy function",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list resource-types called")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	// Subcommands
	destroyCmd.AddCommand(destroyResourceCmd)
	destroyCmd.AddCommand(destroyConnectionCmd)
	destroyCmd.AddCommand(destroyFunctionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
