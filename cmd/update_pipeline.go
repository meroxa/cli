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
	"github.com/meroxa/cli/display"
	"github.com/spf13/cobra"
)

var (
	state string // connector state
)

func UpdatePipelineCmd() *cobra.Command {
	updatePipelineCmd := &cobra.Command{
		Use:     "pipeline <name> --state <pause|resume|restart>",
		Aliases: []string{"pipelines"},
		Short:   "Update pipeline state",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires pipeline name\n\nUsage:\n  meroxa update pipeline <name> --state <pause|resume|restart>")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Pipeline Name
			pipelineName := args[0]

			c, err := client()
			if err != nil {
				return err
			}

			// get pipeline id from name
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			pipeline, err := c.GetPipelineByName(ctx, pipelineName)
			if err != nil {
				return err
			}

			ctx = context.Background()
			ctx, cancel = context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			// call meroxa-go to update pipeline status with name
			if !flagRootOutputJSON {
				fmt.Printf("Updating %s pipeline...\n", pipelineName)
			}

			p, err := c.UpdatePipelineStatus(ctx, pipeline.ID, state)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				display.JSONPrint(p)
			} else {
				fmt.Printf("Pipeline %s successfully updated!\n", p.Name)
			}

			return nil
		},
	}
	updatePipelineCmd.Flags().StringVarP(&state, "state", "", "", "pipeline state")
	updatePipelineCmd.MarkFlagRequired("state")
	return updatePipelineCmd
}