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
	"github.com/spf13/cobra"
)

// RemoveResourceCmd represents the `meroxa remove resource` command
func RemoveResourceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resource <name>",
		Short: "Remove resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires resource name\n\nUsage:\n  meroxa remove resource <name>")
			}
			// Resource Name
			resName := args[0]

			c, err := client()
			if err != nil {
				return err
			}

			// get Resource ID from name
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			res, err := c.GetResourceByName(ctx, resName)
			if err != nil {
				return err
			}

			c, err = client()
			if err != nil {
				return err
			}

			ctx = context.Background()
			ctx, cancel = context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			// TODO: Update meroxa-go to `RemoveResource` to match its implementation
			err = c.DeleteResource(ctx, res.ID)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				utils.JSONPrint(res)
			} else {
				fmt.Printf("Resource %s removed\n", res.Name)
			}
			return nil
		},
	}
}
