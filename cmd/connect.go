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
	"fmt"

	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

// ConnectCmd represents the `meroxa connect` command
func ConnectCmd() *cobra.Command {
	connectCmd := &cobra.Command{
		Use:   "connect --from <resource-name> --to <resource-name>",
		Short: "Connect two resources together",
		Long: `Use the connect command to automatically configure the connectors required to pull data from one resource 
(source) to another (destination).

This command is equivalent to creating two connectors separately, one from the source to Meroxa and another from Meroxa 
to the destination:

meroxa connect --from <resource-name> --to <resource-name> --input <source-input>

or

meroxa create connector --from postgres --input accounts # Creates source connector
meroxa create connector --to redshift --input orders # Creates destination connector
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// config
			cfgString, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			var cfg struct {
				From map[string]string `json:"from"`
				To   map[string]string `json:"to"`
			}
			if cfgString != "" {
				err = json.Unmarshal([]byte(cfgString), &cfg)
				if err != nil {
					return err
				}
			}

			// merge in input
			input, err := cmd.Flags().GetString("input")
			if err != nil {
				return err
			}

			// create connector from source to meroxa
			fmt.Printf("Creating connector from source %s...\n", source)

			// we indicate what type of connector we're creating using its `mx:connectorType` key
			metadata := map[string]string{"mx:connectorType": "source"}

			srcCon, err := createConnector("", source, cfg.From, metadata, input)
			if err != nil {
				return err
			}
			fmt.Printf("Connector from source %s successfully created!\n", source)

			// we use the stream of the source as the input for the destination below
			inputStreams := srcCon.Streams["output"].([]interface{})

			// create connector from meroxa to destination
			fmt.Printf("Creating connector to destination %s...\n", destination)

			metadata["mx:connectorType"] = "destination"
			_, err = createConnector("", destination, cfg.To, metadata, inputStreams[0].(string))
			if err != nil {
				return err
			}
			fmt.Printf("Connector to destination %s successfully created!\n", destination)
			return nil
		},
	}

	connectCmd.Flags().StringVarP(&source, "from", "", "", "source resource name")
	connectCmd.MarkFlagRequired("from")
	connectCmd.Flags().StringVarP(&destination, "to", "", "", "destination resource name")
	connectCmd.MarkFlagRequired("to")
	connectCmd.Flags().StringP("config", "c", "", "connector configuration")
	connectCmd.Flags().String("input", "", "command delimeted list of input streams")

	return connectCmd
}

func createConnector(connectorName string, resourceName string, config map[string]string, metadata map[string]string, input string) (*meroxa.Connector, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}

	// get resource ID from name
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
	defer cancel()

	res, err := c.GetResourceByName(ctx, resourceName)
	if err != nil {
		return nil, err
	}

	// create connector
	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, clientTimeOut)
	defer cancel()

	// merge in input
	if input != "" {
		config["input"] = input
	}

	con, err := c.CreateConnector(ctx, meroxa.CreateConnectorInput{
		Name:          connectorName,
		ResourceID:    res.ID,
		PipelineID:    0,
		Configuration: config,
		Metadata:      metadata,
	})
	if err != nil {
		return nil, err
	}

	return con, nil
}
