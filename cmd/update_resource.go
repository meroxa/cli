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
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

type UpdateResourceClient interface {
	UpdateResource(ctx context.Context, key string, resourceToUpdate meroxa.UpdateResourceInput) (*meroxa.Resource, error)
}

type UpdateResource struct {
	name, metadata, credentials, url string
}

func (ur *UpdateResource) setArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires resource name")
	}

	ur.name = args[0]

	return nil
}

func (ur *UpdateResource) setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ur.url, "url", "u", "", "resource url")
	cmd.Flags().StringVarP(&ur.metadata, "metadata", "m", "", "resource metadata")
	cmd.Flags().StringVarP(&ur.credentials, "credentials", "", "", "resource credentials")
}

func (ur *UpdateResource) execute(ctx context.Context, c UpdateResourceClient) (*meroxa.Resource, error) {
	if ur.url == "" && ur.metadata == "" && ur.credentials == "" {
		return nil, errors.New("requires either `--metadata`, `--url` or `--credentials` to update the resource")
	}

	if !flagRootOutputJSON {
		fmt.Printf("Updating %s resource...\n", ur.name)
	}

	var res meroxa.UpdateResourceInput

	// If url was provided, update it
	if ur.url != "" {
		res.URL = ur.url
	}

	// TODO: Figure out best way to handle creds and metadata
	// Get credentials (expect a JSON string)
	if ur.credentials != "" {
		var creds meroxa.Credentials
		err := json.Unmarshal([]byte(ur.credentials), &creds)
		if err != nil {
			return nil, err
		}

		res.Credentials = &creds
	}

	// If metadata was provided, update it
	if ur.metadata != "" {
		var metadata map[string]string
		err := json.Unmarshal([]byte(ur.metadata), &metadata)
		if err != nil {
			return nil, err
		}

		res.Metadata = metadata
	}

	return c.UpdateResource(ctx, ur.name, res)

	return nil, nil
}

func (ur *UpdateResource) output(res *meroxa.Resource) {
	if flagRootOutputJSON {
		utils.JSONPrint(res)
	} else {
		fmt.Printf("Resource %s successfully updated!\n", ur.name)
	}
}

// UpdateResourceCmd represents the `meroxa update resource` command
func (ur *UpdateResource) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resource <resource-name>",
		Short:   "Update a resource",
		Long:    `Use the update command to update various Meroxa resources.`,
		Aliases: []string{"resources"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ur.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()

			if err != nil {
				return err
			}

			resource, err := ur.execute(ctx, c)
			if err != nil {
				return err
			}

			ur.output(resource)
			return nil
		},
	}

	ur.setFlags(cmd)

	return cmd
}
