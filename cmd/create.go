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

		if !flagRootOutputJson {
			fmt.Printf("Creating %s Resource...\n", resType)
		}

		res, err := c.CreateResource(ctx, &r)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if flagRootOutputJson {
			jsonPrint(res)
		} else {
			fmt.Println("Resource successfully created!")
			prettyPrint("resource", res)
		}
	},
}

var createConnectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "create connection",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Resource Name
		resName := args[0]

		cfgString, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		cfg := &Config{}
		if cfgString != "" {
			err = json.Unmarshal([]byte(cfgString), cfg)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}

		// merge in input
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		err = cfg.Set("input", input)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if !flagRootOutputJson {
			fmt.Println("Creating connection...")
		}

		con, err := createConnection(resName, cfg, input)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if flagRootOutputJson {
			jsonPrint(con)
		} else {
			fmt.Println("Connection successfully created!")
			prettyPrint("connector", con)
		}
	},
}

var createFunctionCmd = &cobra.Command{
	Use:   "function",
	Short: "create function",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create function called - Not Implemented")
	},
}

var createPipelineCmd = &cobra.Command{
	Use:   "pipeline <name>",
	Short: "create pipeline",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pipelineName := args[0]

		c, err := client()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		p := &meroxa.Pipeline{
			Name: pipelineName,
		}

		// Process metadata
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
			p.Metadata = metadata
		}

		if !flagRootOutputJson {
			fmt.Println("Creating Pipeline...")
		}

		res, err := c.CreatePipeline(ctx, p)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if flagRootOutputJson {
			jsonPrint(res)
		} else {
			fmt.Println("Pipeline successfully created!")
			prettyPrint("pipeline", res)
		}
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
	createConnectionCmd.Flags().String("input", "", "command delimeted list of input streams")
	createConnectionCmd.MarkFlagRequired("input")

	createCmd.AddCommand(createFunctionCmd)

	createCmd.AddCommand(createPipelineCmd)
	createPipelineCmd.Flags().StringP("metadata", "m", "", "pipeline metadata")
}

func createConnection(resourceName string, config *Config, input string) (*meroxa.Connector, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}

	// get resource ID from name
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.GetResourceByName(ctx, resourceName)
	if err != nil {
		return nil, err
	}

	// create connection
	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	//var cfg Config
	cfg := Config{}
	if config != nil {
		cfg.Merge(*config)
	}

	// merge in input
	if input != "" {
		err = cfg.Set("input", input)
		if err != nil {
			return nil, err
		}
	}

	con, err := c.CreateConnection(ctx, res.ID, cfg)
	if err != nil {
		return nil, err
	}

	return con, nil
}
