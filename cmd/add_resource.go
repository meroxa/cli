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

type AddResourceClient interface {
	CreateResource(ctx context.Context, resource *meroxa.CreateResourceInput) (*meroxa.Resource, error)
}

type AddResource struct {
	name, rType, url, metadata, credentials string
}

func (ar *AddResource) setArgs(args []string) error {
	if len(args) > 0 {
		ar.name = args[0]
	}

	return nil
}

func (ar *AddResource) setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ar.rType, "type", "", "", "resource type")
	cmd.MarkFlagRequired("type")

	cmd.Flags().StringVarP(&ar.url, "url", "u", "", "resource url")
	cmd.MarkFlagRequired("url")

	cmd.Flags().StringVarP(&ar.credentials, "credentials", "", "", "resource credentials")
	cmd.Flags().StringVarP(&ar.metadata, "metadata", "m", "", "resource metadata")
}

func (ar *AddResource) execute(ctx context.Context, c AddResourceClient, res meroxa.CreateResourceInput) (*meroxa.Resource, error) {
	if !flagRootOutputJSON {
		fmt.Printf("Adding %s resource...\n", res.Type)
	}

	var err error

	// TODO: Figure out best way to handle creds and metadata
	// Get credentials (expect a JSON string)
	if ar.credentials != "" {
		var creds meroxa.Credentials
		err = json.Unmarshal([]byte(ar.credentials), &creds)
		if err != nil {
			return nil, err
		}

		res.Credentials = &creds
	}

	if ar.metadata != "" {
		var metadata map[string]string
		err = json.Unmarshal([]byte(ar.metadata), &metadata)
		if err != nil {
			return nil, err
		}

		res.Metadata = metadata
	}

	resource, err := c.CreateResource(ctx, &res)
	return resource, err
}

func (ar *AddResource) output(res *meroxa.Resource) {
	if flagRootOutputJSON {
		utils.JSONPrint(res)
	} else {
		fmt.Printf("%s resource with name %s successfully added!\n", res.Type, res.Name)
	}
}

// AddResourceCmd represents the `meroxa add resource` command
func (ar *AddResource) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource [NAME] --type TYPE",
		Short: "Add a resource to your Meroxa resource catalog",
		Long:  `Use the add command to add resources to your Meroxa resource catalog.`,
		Example: "\n" +
			"meroxa add resource store --type postgres -u $DATABASE_URL\n" +
			"meroxa add resource datalake --type s3 -u \"s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos\"\n" +
			"meroxa add resource warehouse --type redshift -u $REDSHIFT_URL\n" +
			"meroxa add resource slack --type url -u $WEBHOOK_URL\n",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ar.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()

			if err != nil {
				return err
			}

			ri := meroxa.CreateResourceInput{
				Type:     ar.rType,
				Name:     ar.name,
				URL:      ar.url,
				Metadata: nil,
			}

			res, err := ar.execute(ctx, c, ri)

			if err != nil {
				return err
			}

			ar.output(res)

			return nil
		},
	}

	ar.setFlags(cmd)

	return cmd
}
