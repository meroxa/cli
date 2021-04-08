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

type RemovePipeline struct {
	name      string
	removeCmd *Remove
}

// RemovePipelineClient represents the interface for meroxa client
type RemovePipelineClient interface {
	GetPipelineByName(ctx context.Context, name string) (*meroxa.Pipeline, error)
	DeletePipeline(ctx context.Context, id int) error
}

func (rp *RemovePipeline) setArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires pipeline name\n\nUsage:\n  meroxa remove pipeline NAME")
	}
	// endpoint name
	rp.name = args[0]

	rp.removeCmd.componentType = "pipeline"
	rp.removeCmd.confirmableName = rp.name

	return nil
}

func (rp *RemovePipeline) execute(ctx context.Context, c RemovePipelineClient) (*meroxa.Pipeline, error) {
	p, err := c.GetPipelineByName(ctx, rp.name)
	if err != nil {
		return nil, err
	}

	return p, c.DeletePipeline(ctx, p.ID)
}

func (rp *RemovePipeline) output(p *meroxa.Pipeline) {
	if flagRootOutputJSON {
		utils.JSONPrint(p)
	} else {
		fmt.Printf("pipeline %s successfully removed\n", p.Name)
	}
}

// command represents the `meroxa remove pipeline` command
func (rp *RemovePipeline) command() *cobra.Command {
	return &cobra.Command{
		Use:   "pipeline NAME",
		Short: "Remove pipeline",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return rp.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()
			if err != nil {
				return err
			}

			p, err := rp.execute(ctx, c)

			rp.output(p)

			return nil
		},
	}
}
