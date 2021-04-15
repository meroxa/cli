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
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

type ListConnectorsClient interface {
	ListConnectors(ctx context.Context) ([]*meroxa.Connector, error)
	ListPipelineConnectors(ctx context.Context, pipelineID int) ([]*meroxa.Connector, error)
	GetPipelineByName(ctx context.Context, name string) (*meroxa.Pipeline, error)
}

type ListConnectors struct {
	pipeline string
}

func (lc *ListConnectors) execute(ctx context.Context, c ListConnectorsClient) ([]*meroxa.Connector, error) {
	if lc.pipeline != "" {
		p, err := c.GetPipelineByName(ctx, lc.pipeline)

		if err != nil {
			return nil, err
		}

		return c.ListPipelineConnectors(ctx, p.ID)
	}

	return c.ListConnectors(ctx)
}

func (lc *ListConnectors) setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&lc.pipeline, "pipeline", "", "", "filter connectors by pipeline name")
}

func (lc *ListConnectors) output(connectors []*meroxa.Connector) {
	if FlagRootOutputJSON {
		utils.JSONPrint(connectors)
	} else {
		utils.PrintConnectorsTable(connectors)
	}
}

// ListConnectorsCmd represents the `meroxa list connectors` command
func (lc *ListConnectors) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connectors",
		Short:   "List connectors",
		Aliases: []string{"connector"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := global.NewClient()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, ClientTimeOut)
			defer cancel()

			connectors, err := lc.execute(ctx, c)
			if err != nil {
				return err
			}

			lc.output(connectors)

			return nil
		},
	}

	lc.setFlags(cmd)
	return cmd
}
