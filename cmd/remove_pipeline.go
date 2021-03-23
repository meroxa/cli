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

// RemovePipelineCmd represents the `meroxa remove pipeline` command
func RemovePipelineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pipeline <name>",
		Short: "Remove pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires pipeline name\n\nUsage:\n  meroxa remove pipeline <name>")
			}

			// Pipeline Name
			pipelineName := args[0]

			c, err := client()
			if err != nil {
				return err
			}

			// get Pipeline ID from name
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			pipeline, err := c.GetPipelineByName(ctx, pipelineName)
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

			err = c.DeletePipeline(ctx, pipeline.ID)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				utils.JSONPrint(pipeline)
			} else {
				fmt.Printf("Pipeline %s removed\n", pipeline.Name)
			}

			return nil
		},
	}
}
