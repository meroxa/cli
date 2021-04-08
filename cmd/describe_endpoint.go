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
	"fmt"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

// DescribeEndpointCmd represents the `meroxa describe endpoint` command
func DescribeEndpointCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "endpoint NAME",
		Aliases: []string{"endpoints"},
		Short:   "Describe endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires endpoint name\n\nUsage:\n  meroxa describe endpoint NAME [flags]")
			}
			name := args[0]

			c, err := client()
			if err != nil {
				return err
			}
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			end, err := c.GetEndpoint(ctx, name)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				utils.JSONPrint(end)
			} else {
				utils.PrintEndpointsTable([]meroxa.Endpoint{*end})
			}
			return nil

		},
	}
}
