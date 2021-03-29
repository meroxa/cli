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

	"github.com/meroxa/cli/utils"

	"github.com/spf13/cobra"
)

// ListResourcesCmd represents the `meroxa list resources` command
func ListResourcesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "resources",
		Short:   "List resources",
		Aliases: []string{"resource"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			rr, err := c.ListResources(ctx)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				utils.JSONPrint(rr)
			} else {
				utils.PrintResourcesTable(rr)
			}
			return nil
		},
	}
}
