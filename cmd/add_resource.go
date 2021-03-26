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
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

type addResourceCmd struct {

}

type AddResourceClient interface {
	CreateResource(ctx context.Context, resource *meroxa.CreateResourceInput) (*meroxa.Resource, error)
}

func addResourceArgs (args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	return "", nil
}

func addResourceFlags (cmd *cobra.Command) *cobra.Command{
	cmd.Flags().StringVarP(&resType, "type", "", "", "resource type")
	cmd.MarkFlagRequired("type")

	cmd.Flags().StringVarP(&resURL, "url", "u", "", "resource url")
	cmd.MarkFlagRequired("url")

	cmd.Flags().StringVarP(&resCredentials, "credentials", "", "", "resource credentials")
	cmd.Flags().StringVarP(&resMetadata, "metadata", "m", "", "resource metadata")

	return cmd
}

func addResource (c AddResourceClient, rType, rName, rURL string) (*meroxa.Resource, error) {
	var err error

	r := meroxa.CreateResourceInput{
		Type:     rType,
		Name:     rName,
		URL:      rURL,
		Metadata: nil,
	}

	// TODO: Figure out best way to handle creds and metadata
	// Get credentials (expect a JSON string)
	if resCredentials != "" {
		var creds meroxa.Credentials
		err = json.Unmarshal([]byte(resCredentials), &creds)
		if err != nil {
			return nil, err
		}

		r.Credentials = &creds
	}

	if resMetadata != "" {
		var metadata map[string]string
		err = json.Unmarshal([]byte(resMetadata), &metadata)
		if err != nil {
			return nil, err
		}

		r.Metadata = metadata
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
	defer cancel()

	if !flagRootOutputJSON {
		fmt.Printf("Adding %s resource (%s)...\n", resName, resType)
	}

	return c.CreateResource(ctx, &r)
}

func addResourceOutput(r *meroxa.Resource) {
	if flagRootOutputJSON {
		utils.JSONPrint(r)
	} else {
		fmt.Printf("Resource %s successfully added!\n", r.Name)
	}
}

// AddResourceCmd represents the `meroxa add resource` command
func AddResourceCmd() *cobra.Command {
	addResourceCmd := &cobra.Command{
		Use:   "resource <resource-name> --type <resource-type>",
		Short: "Add a resource to your Meroxa resource catalog",
		Long:  `Use the add command to add resources to your Meroxa resource catalog.`,
		Example: "\n" +
			"meroxa add resource store --type postgres -u $DATABASE_URL\n" +
			"meroxa add resource datalake --type s3 -u \"s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos\"\n" +
			"meroxa add resource warehouse --type redshift -u $REDSHIFT_URL\n" +
			"meroxa add resource slack --type url -u $WEBHOOK_URL\n",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error

			resName, err = addResourceArgs(args)

			if err != nil {
				cmd.PrintErr(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client()

			if err != nil {
				return err
			}

			res, err := addResource(c, resType, resName, resURL)

			if err != nil {
				return err
			}

			addResourceOutput(res)
			return nil
		},
	}

	addResourceCmd = addResourceFlags(addResourceCmd)

	return addResourceCmd
}