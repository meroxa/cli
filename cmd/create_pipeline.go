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

// CreatePipelineCmd represents the `meroxa create pipeline` command
func CreatePipelineCmd() *cobra.Command {
	createPipelineCmd := &cobra.Command{
		Use:   "pipeline NAME",
		Short: "Create a pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a pipeline name\n\nUsage:\n  meroxa create pipeline NAME [flags]")
			}
			pipelineName := args[0]

			c, err := client()
			if err != nil {
				return err
			}
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			p := &meroxa.Pipeline{
				Name: pipelineName,
			}

			// Process metadata
			metadataString, err := cmd.Flags().GetString("metadata")
			if err != nil {
				return err
			}
			if metadataString != "" {
				var metadata map[string]string
				err = json.Unmarshal([]byte(metadataString), &metadata)
				if err != nil {
					return err
				}
				p.Metadata = metadata
			}

			if !flagRootOutputJSON {
				fmt.Println("Creating Pipeline...")
			}

			res, err := c.CreatePipeline(ctx, p)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				utils.JSONPrint(res)
			} else {
				fmt.Printf("Pipeline %s successfully created!\n", p.Name)
			}
			return nil
		},
	}

	createPipelineCmd.Flags().StringP("metadata", "m", "", "pipeline metadata")
	return createPipelineCmd
}
