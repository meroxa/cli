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
	"encoding/json"
	"fmt"
	"time"

	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create meroxa pipeline components",
	Long: `use the create command to create various Meroxa pipeline components
including Resources, Connections and Functions.`,
}

var createResourceCmd = &cobra.Command{
	Use:   "resource <resource-type>",
	Short: "create resource",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Resource Type
		resType := args[0]

		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		// Assemble resource struct from config
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error: ", err)
		}

		u, err := cmd.Flags().GetString("url")
		if err != nil {
			fmt.Println("Error: ", err)
		}

		r := meroxa.Resource{
			Kind:          resType,
			Name:          name,
			URL:           u,
			Configuration: nil,
			Metadata:      nil,
		}

		// TODO: Figure out best way to handle creds, config and metadata
		// Get credentials (expect a JSON string)
		credsString, err := cmd.Flags().GetString("credentials")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if credsString != "" {
			var creds meroxa.Credentials
			err = json.Unmarshal([]byte(credsString), &creds)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			r.Credentials = &creds
		}

		metadataString, err := cmd.Flags().GetString("metadata")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if metadataString != "" {
			var metadata map[string]string
			err = json.Unmarshal([]byte(metadataString), &metadata)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			r.Metadata = metadata
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		fmt.Printf("Creating %s Resource...\n", resType)

		res, err := c.CreateResource(ctx, &r)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Println("Resource successfully created!")
		prettyPrint("resource", res)
	},
}

var createConnectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "create connection",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Resource Name
		resName := args[0]

		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// get resource ID from name
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		res, err := c.GetResourceByName(ctx, resName)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// create connection
		ctx = context.Background()
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		cfgString, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		var cfg map[string]string
		err = json.Unmarshal([]byte(cfgString), &cfg)
		if err != nil {
			fmt.Println("1Error: ", err)
			return
		}

		fmt.Println("Creating connection...")
		con, err := c.CreateConnection(ctx, res.ID, cfg)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		fmt.Println("Connection successfully created!")
		prettyPrint("connector", con)
	},
}

var createFunctionCmd = &cobra.Command{
	Use:   "function",
	Short: "create function",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list resource-types called")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Subcommands
	createCmd.AddCommand(createResourceCmd)
	createResourceCmd.Flags().StringP("name", "n", "", "resource name")
	createResourceCmd.Flags().StringP("url", "u", "", "resource url")
	createResourceCmd.Flags().String("credentials", "", "resource credentials")
	createResourceCmd.Flags().StringP("config", "c", "", "resource configuration")
	createResourceCmd.Flags().StringP("metadata", "m", "", "resource metadata")

	createCmd.AddCommand(createConnectionCmd)
	createConnectionCmd.Flags().StringP("config", "c", "", "connection configuration")
	createCmd.AddCommand(createFunctionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
