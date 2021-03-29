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
	"errors"
	"fmt"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

type RemoveResource struct {
	name string
}

type RemoveResourceClient interface {
	GetResourceByName(ctx context.Context, name string) (*meroxa.Resource, error)
	DeleteResource(ctx context.Context, id int) error
}

var removeResourceCmd RemoveResource

func (RemoveResource) setArgs (args []string) error {
	if len(args) < 1 {
		return errors.New("requires resource name\n\nUsage:\n  meroxa remove resource <name>")
	}
	// Resource Name
	removeResourceCmd.name = args[0]
	return nil
}

func (RemoveResource) execute (ctx context.Context, c RemoveResourceClient) (*meroxa.Resource, error) {
	// get Resource ID from name
	res, err := c.GetResourceByName(ctx, resName)
	if err != nil {
		return nil, err
	}

	c, err = client()

	if err != nil {
		return nil, err
	}

	if !removeCmd.force {
		return nil, errors.New("removing resource not confirmed")
	}

	return res, c.DeleteResource(ctx, res.ID)
}

func (RemoveResource) output(r *meroxa.Resource) {
	if flagRootOutputJSON {
		utils.JSONPrint(r)
	} else {
		fmt.Printf("Resource %s removed\n", r.Name)
	}
}

// RemoveResource represents the `meroxa remove resource` command
func (RemoveResource) command() *cobra.Command {
	return &cobra.Command{
		Use:   "resource <name>",
		Short: "Remove resource",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return removeResourceCmd.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()

			if err != nil {
				return err
			}

			var r *meroxa.Resource
			r, err = removeResourceCmd.execute(ctx, c)

			if err != nil {
				return err
			}

			removeResourceCmd.output(r)

			return nil
		},
	}
}
