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

type addResourceClient interface {
	CreateResource(ctx context.Context, resource *meroxa.CreateResourceInput) (*meroxa.Resource, error)
}

type AddResource struct{
	name, rType, url, metadata, credentials string
}

func (AddResource) checkArgs (args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	return "", nil
}

func (c AddResource) setFlags (cmd *cobra.Command) {
	cmd.Flags().StringVarP(&c.rType, "type", "", "", "resource type")
	cmd.MarkFlagRequired("type")

	cmd.Flags().StringVarP(&c.url, "url", "u", "", "resource url")
	cmd.MarkFlagRequired("url")

	cmd.Flags().StringVarP(&c.credentials, "credentials", "", "", "resource credentials")
	cmd.Flags().StringVarP(&c.metadata, "metadata", "m", "", "resource metadata")
}

func (c AddResource) execute (ctx context.Context, client addResourceClient, r meroxa.CreateResourceInput) (*meroxa.Resource, error) {
	var err error

	// TODO: Figure out best way to handle creds and metadata
	// Get credentials (expect a JSON string)
	if resCredentials != "" {
		var creds meroxa.Credentials
		err = json.Unmarshal([]byte(c.credentials), &creds)
		if err != nil {
			return nil, err
		}

		r.Credentials = &creds
	}

	if c.metadata != "" {
		var metadata map[string]string
		err = json.Unmarshal([]byte(c.metadata), &metadata)
		if err != nil {
			return nil, err
		}

		r.Metadata = metadata
	}

	if !flagRootOutputJSON {
		fmt.Printf("Adding %s resource...\n", r.Type)
	}

	return client.CreateResource(ctx, &r)
}

func (AddResource) output(r *meroxa.Resource) {
	if flagRootOutputJSON {
		utils.JSONPrint(r)
	} else {
		fmt.Printf("Resource %s successfully added!\n", r.Name)
	}
}

// AddResourceCmd represents the `meroxa add resource` command
func (c AddResource) command() *cobra.Command {
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

			c.name, err = c.checkArgs(args)

			if err != nil {
				cmd.PrintErr(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			cl, err := client()

			if err != nil {
				return err
			}

			r := meroxa.CreateResourceInput{
				Type:     c.rType,
				Name:     c.name,
				URL:      c.url,
				Metadata: nil,
			}

			res, err := c.execute(ctx, cl, r)

			if err != nil {
				return err
			}

			c.output(res)

			return nil
		},
	}

	c.setFlags(addResourceCmd)

	return addResourceCmd
}