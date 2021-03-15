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
	"errors"
	"fmt"
	"time"

	"github.com/meroxa/cli/display"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

var updateResourceCmd = &cobra.Command{
	Use:     "resource <resource-name>",
	Short:   "Update a resource",
	Long:    `Use the update command to update various Meroxa resources.`,
	Aliases: []string{"resources"},
	// TODO: Change the design so a new name for the resource could be set
	// meroxa update resource <old-resource-name> --name <new-resource-name>
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || (resURL == "" && resMetadata == "" && resCredentials == "") {
			return errors.New("requires a resource name and either `--metadata`, `--url` or `--credentials` to update the resource \n\nUsage:\n  meroxa update resource <resource-name> [--url <url> | --metadata <metadata> | --credentials <credentials>]")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Resource Name
		resName = args[0]
		c, err := client()

		if err != nil {
			return err
		}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var res meroxa.UpdateResourceInput

		// If url was provided, update it
		if resURL != "" {
			res.URL = resURL
		}

		// TODO: Figure out best way to handle creds and metadata
		// Get credentials (expect a JSON string)
		if resCredentials != "" {
			var creds meroxa.Credentials
			err = json.Unmarshal([]byte(resCredentials), &creds)
			if err != nil {
				return err
			}

			res.Credentials = &creds
		}

		// If metadata was provided, update it
		if resMetadata != "" {
			var metadata map[string]string
			err = json.Unmarshal([]byte(resMetadata), &metadata)
			if err != nil {
				return err
			}

			res.Metadata = metadata
		}

		// call meroxa-go to update resource
		if !flagRootOutputJSON {
			fmt.Printf("Updating %s resource...\n", resName)
		}

		resource, err := c.UpdateResource(ctx, resName, res)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(resource)
		} else {
			fmt.Printf("Resource %s successfully updated!\n", resName)
		}

		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateResourceCmd)

	updateResourceCmd.Flags().StringVarP(&resURL, "url", "u", "", "resource url")
	updateResourceCmd.Flags().StringVarP(&resMetadata, "metadata", "m", "", "resource metadata")
	updateResourceCmd.Flags().StringVarP(&resCredentials, "credentials", "", "", "resource credentials")
}
