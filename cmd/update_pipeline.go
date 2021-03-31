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

type UpdatePipelineClient interface {
	GetPipelineByName(ctx context.Context, name string) (*meroxa.Pipeline, error)
	// TODO: Try to unify UpdatePipelineStatus with UpdatePipeline in meroxa-go
	UpdatePipelineStatus(ctx context.Context, pipelineID int, state string) (*meroxa.Pipeline, error)
	UpdatePipeline(ctx context.Context, key string, pipeline meroxa.UpdatePipelineInput) (*meroxa.Pipeline, error)
}

type UpdatePipeline struct {
	name, newName, metadata, state string
}

func (up *UpdatePipeline) setArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires pipeline name")
	}

	up.name = args[0]

	// TODO: Depending on what's been updated, set values here so execute knows what to update
	return nil
}

func (up *UpdatePipeline) setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&up.state, "state", "", "", "new pipeline state")
	cmd.Flags().StringVarP(&up.newName, "name", "", "", "new pipeline name")
	cmd.Flags().StringVarP(&up.metadata, "metadata", "m", "", "new pipeline metadata")
}

func (up *UpdatePipeline) execute(ctx context.Context, c UpdatePipelineClient) (*meroxa.Pipeline, error) {
	if up.newName == "" && up.metadata == "" && up.state == "" {
		return nil, errors.New("requires either --name, --state or --metadata")
	}

	if !flagRootOutputJSON {
		fmt.Printf("Updating %s pipeline...\n", up.name)
	}

	// get pipeline id from name
	p, err := c.GetPipelineByName(ctx, up.name)

	if err != nil {
		return p, err
	}

	// call meroxa-go to update pipeline status with name
	p, err = c.UpdatePipelineStatus(ctx, p.ID, up.state)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (up *UpdatePipeline) output(p *meroxa.Pipeline) {
	if flagRootOutputJSON {
		utils.JSONPrint(p)
	} else {
		fmt.Printf("pipeline %s successfully updated!\n", p.Name)
	}
}

// command represents the `meroxa update pipeline` command
func (up *UpdatePipeline) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pipeline <name>",
		Aliases: []string{"pipelines"},
		Short:   "Update pipeline state",
		Example: "\n" +
			"meroxa update pipeline old-name --name new-name\n" +
			"meroxa update pipeline pipeline-name --state pause\n" +
			"meroxa update pipeline pipeline-name --metadata '{\"key\":\"value\"}'",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return up.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()
			if err != nil {
				return err
			}

			p, err := up.execute(ctx, c)

			if err != nil {
				return err
			}

			up.output(p)
			return nil
		},
	}

	up.setFlags(cmd)
	return cmd
}
