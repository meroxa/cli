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
	"context"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/utils"

	"github.com/spf13/cobra"
)

// ListPipelinesCmd represents the `meroxa list pipelines` command
func ListPipelinesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "pipelines",
		Short:   "List pipelines",
		Aliases: []string{"pipeline"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := global.NewClient()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, ClientTimeOut)
			defer cancel()

			rr, err := c.ListPipelines(ctx)
			if err != nil {
				return err
			}

			if FlagRootOutputJSON {
				utils.JSONPrint(rr)
			} else {
				utils.PrintPipelinesTable(rr)
			}
			return nil
		},
	}
}
