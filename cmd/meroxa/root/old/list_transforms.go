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

package old

import (
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/utils"
	"github.com/spf13/cobra"
)

// ListTransformsCmd represents the `meroxa list transforms` command.
func ListTransformsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "transforms",
		Short:   "List transforms",
		Aliases: []string{"transform"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := global.NewClient()
			if err != nil {
				return err
			}

			rr, err := c.ListTransforms(cmd.Context())
			if err != nil {
				return err
			}

			if FlagRootOutputJSON {
				utils.JSONPrint(rr)
			} else {
				utils.PrintTransformsTable(rr)
			}
			return nil
		},
	}
}
