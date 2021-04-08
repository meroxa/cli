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

	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

type Connect struct {
	input, config, source, destination, pipelineName string
}

func (conn *Connect) setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&conn.source, "from", "", "", "source resource name")
	cmd.MarkFlagRequired("from")
	cmd.Flags().StringVarP(&conn.destination, "to", "", "", "destination resource name")
	cmd.MarkFlagRequired("to")
	cmd.Flags().StringVarP(&conn.config, "config", "c", "", "connector configuration")
	cmd.Flags().StringVarP(&conn.input, "input", "", "", "command delimeted list of input streams")
	cmd.Flags().StringVarP(&conn.pipelineName, "pipeline", "", "", "pipeline name to attach connectors to")
}

func (conn *Connect) execute(ctx context.Context, c CreateConnectorClient) (*meroxa.Connector, *meroxa.Connector, error) {
	srcCc := &CreateConnector{
		input:        conn.input,
		config:       conn.config,
		source:       conn.source,
		pipelineName: conn.pipelineName,
	}

	srcCon, err := srcCc.execute(ctx, c)
	if err != nil {
		return nil, nil, err
	}
	srcCc.output(srcCon)

	// we use the stream of the source as the input for the destination below
	inputStreams := srcCon.Streams["output"].([]interface{})

	destCc := &CreateConnector{
		input:        inputStreams[0].(string),
		config:       conn.config,
		destination:  conn.destination,
		pipelineName: conn.pipelineName,
	}

	destCon, err := destCc.execute(ctx, c)
	if err != nil {
		return nil, nil, err
	}
	destCc.output(destCon)

	return srcCon, destCon, nil
}

// command returns the cobra Command for `connect`
func (conn *Connect) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect --from RESOURCE-NAME --to RESOURCE-NAME",
		Short: "Connect two resources together",
		Long: `Use the connect command to automatically configure the connectors required to pull data from one resource 
(source) to another (destination).

This command is equivalent to creating two connectors separately, one from the source to Meroxa and another from Meroxa 
to the destination:

meroxa connect --from RESOURCE-NAME --to RESOURCE-NAME --input SOURCE-INPUT

or

meroxa create connector --from postgres --input accounts # Creates source connector
meroxa create connector --to redshift --input orders # Creates destination connector
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()
			if err != nil {
				return err
			}

			_, _, err = conn.execute(ctx, c)
			if err != nil {
				return err
			}

			return nil
		},
	}

	conn.setFlags(cmd)

	return cmd
}
