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
	"github.com/meroxa/cli/display"

	"time"

	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

const clientTimeOut = 5 * time.Second

var resName, resType, resURL, resCredentials, resMetadata string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a resource to your Meroxa resource catalog",
}

var addResourceCmd = &cobra.Command{
	Use:   "resource <resource-name> --type <resource-type>",
	Short: "Add a resource to your Meroxa resource catalog",
	Long:  `Use the add command to add resources to your Meroxa resource catalog.`,
	Example: "\n" +
		"meroxa add resource store --type postgres -u $DATABASE_URL\n" +
		"meroxa add resource datalake --type s3 -u \"s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos\"\n" +
		"meroxa add resource warehouse --type redshift -u $REDSHIFT_URL\n" +
		"meroxa add resource slack --type url -u $WEBHOOK_URL\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires resource name\n\nUsage:\n  meroxa add resource <resource-name> [flags]")
		}

		// If using newer version of the CLI
		if resType != "" {
			resName = args[0]
			fmt.Printf("Using new version, sets resName as %s with resType %s", args[0], resType)
		} else {
			resType = args[0]
			fmt.Printf("Using old version, sets resType as %s with resName %s", args[0], resName)
		}
		c, err := client()

		if err != nil {
			return err
		}

		r := meroxa.Resource{
			Kind:     resType,
			Name:     resName,
			URL:      resURL,
			Metadata: nil,
		}

		// TODO: Figure out best way to handle creds and metadata
		// Get credentials (expect a JSON string)
		if resCredentials != "" {
			var creds meroxa.Credentials
			err = json.Unmarshal([]byte(resCredentials), &creds)
			if err != nil {
				return err
			}

			r.Credentials = &creds
		}

		if resMetadata != "" {
			var metadata map[string]string
			err = json.Unmarshal([]byte(resMetadata), &metadata)
			if err != nil {
				return err
			}

			r.Metadata = metadata
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
		defer cancel()

		if !flagRootOutputJSON {
			display.PrettyPrint("Adding %s resource...\n", resType)
		}

		res, err := c.CreateResource(ctx, &r)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(res)
		} else {
			display.PrettyPrint("Resource %s successfully added!\n", res.Name)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addResourceCmd)

	addResourceCmd.Flags().StringVarP(&resType, "type", "", "", "resource type")
	// TODO: Remove this flag as optional once we release a newer version
	// addResourceCmd.MarkFlagRequired("type")

	// TODO: Remove this flag altogether once we release a newer version
	addResourceCmd.Flags().StringVarP(&resName, "name", "", "", "resource name")
	addResourceCmd.Flags().MarkHidden("name")

	addResourceCmd.Flags().StringVarP(&resURL, "url", "u", "", "resource url")
	addResourceCmd.MarkFlagRequired("url")

	addResourceCmd.Flags().StringVarP(&resCredentials, "credentials", "", "", "resource credentials")
	addResourceCmd.Flags().StringVarP(&resMetadata, "metadata", "m", "", "resource metadata")
}
