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
	"encoding/json"
	"fmt"
	"time"

	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

func addResource(resType string, cmd *cobra.Command) {
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
		Kind: resType,
		Name: name,
		URL:  u,
		// We're not doing anything with `config` in the CLI.
		// Maybe deprecate this altogether in the client.
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
	const FIVE = 5
	ctx, cancel := context.WithTimeout(ctx, FIVE*time.Second)
	defer cancel()

	if !flagRootOutputJSON {
		fmt.Printf("Creating %s Resource...\n", resType)
	}

	res, err := c.CreateResource(ctx, &r)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if flagRootOutputJSON {
		jsonPrint(res)
	} else {
		fmt.Println("Resource successfully created!")
		prettyPrint("resource", res)
	}
}

// addResourceCmd represents the add command
var addResourceCmd = &cobra.Command{
	Use:   "add resource <resource-type>",
	Short: "Add a resource to your account",
	Run: func(cmd *cobra.Command, args []string) {
		addResource(args[0], cmd)
	},
}

func init() {
	rootCmd.AddCommand(addResourceCmd)

	addResourceCmd.Flags().StringP("name", "n", "", "resource name")
	addResourceCmd.Flags().StringP("url", "u", "", "resource url")
	addResourceCmd.Flags().String("credentials", "", "resource credentials")
	addResourceCmd.Flags().StringP("config", "c", "", "resource configuration")
	addResourceCmd.Flags().StringP("metadata", "m", "", "resource metadata")
}
