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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
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

		cfg := struct {
			From *Config `json:"from"`
			To   *Config `json:"to"`
		}{
			From: &Config{},
			To:   &Config{},
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
		fmt.Printf("Creating connector from source %s...", source)

		// we indicate what type of connector we're creating using its `mx:connectorType` key
		metadata := map[string]string{"mx:connectorType": ""}
		metadata["mx:connectorType"] = "source"

		srcCon, err := createConnector("", source, cfg.From, metadata, input)
		if err != nil {
			return err
		}
		fmt.Printf("Connector from source %s successfully created!", source)

		// we use the stream of the source as the input for the destination below
		inputStreams := srcCon.Streams["output"].([]interface{})

		// create connector from meroxa to destination
		fmt.Printf("Creating connector to destination %s...", destination)

		metadata["mx:connectorType"] = "destination"
		_, err = createConnector("", destination, cfg.To, metadata, inputStreams[0].(string))
		if err != nil {
			return err
		}
		fmt.Printf("Connector to destination %s successfully created!", destination)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(connectCmd)

	// Flags
	connectCmd.Flags().StringVarP(&source, "from", "", "", "source resource name")
	connectCmd.MarkFlagRequired("from")
	connectCmd.Flags().StringVarP(&destination, "to", "", "", "destination resource name")
	connectCmd.MarkFlagRequired("to")
	connectCmd.Flags().StringP("config", "c", "", "connector configuration")
	connectCmd.Flags().String("input", "", "command delimeted list of input streams")
}
