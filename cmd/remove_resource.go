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
	"os"
)

type RemoveResource struct {
	name      string
	removeCmd *Remove
}

type RemoveResourceClient interface {
	GetResourceByName(ctx context.Context, name string) (*meroxa.Resource, error)
	DeleteResource(ctx context.Context, id int) error
}

func (rr *RemoveResource) setArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires resource name\n\nUsage:\n  meroxa remove resource <name>")
	}
	// Resource Name
	rr.name = args[0]
	return nil
}

func (rr *RemoveResource) execute(ctx context.Context, c RemoveResourceClient) (*meroxa.Resource, error) {
	if !flagRootOutputJSON {
		fmt.Printf("Removing resource %s...\n", rr.name)
	}

	// get Resource ID from name
	res, err := c.GetResourceByName(ctx, rr.name)
	if err != nil {
		return nil, err
	}

	canRemove := rr.removeCmd.confirmRemove(os.Stdin, rr.name)

	if canRemove {
		return res, c.DeleteResource(ctx, res.ID)
	}

	return res, errors.New("removing resource not confirmed")
}

func (rr *RemoveResource) output(r *meroxa.Resource) {
	if flagRootOutputJSON {
		utils.JSONPrint(r)
	} else {
		fmt.Printf("resource %s successfully removed\n", r.Name)
	}
}

// RemoveResource represents the `meroxa remove resource` command
func (rr *RemoveResource) command() *cobra.Command {
	return &cobra.Command{
		Use:   "resource <name>",
		Short: "Remove resource",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return rr.setArgs(args)
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
			r, err = rr.execute(ctx, c)

			if err != nil {
				return err
			}

			rr.output(r)

			return nil
		},
	}
}
