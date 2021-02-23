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
	Long: `Use the connect commands to automatically configure the connectors
required to pull data from one resource (the source) to another
(the target).

This is essentially a shortcut for creating a connector from the
source to Meroxa and creating a connector from Meroxa to the target`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// source name
		sourceName, err := cmd.Flags().GetString("from")
		if err != nil {
			return err
		}

		// target name
		targetName, err := cmd.Flags().GetString("to")
		if err != nil {
			return err
		}

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
		fmt.Println("Creating connector from source...")
		srcCon, err := createConnector("", sourceName, cfg.From, nil, input)
		if err != nil {
			return err
		}
		fmt.Println("Connector successfully created!")

		inputStreams := srcCon.Streams["output"].([]interface{})

		fmt.Println("Creating connector to target...")
		_, err = createConnector("", targetName, cfg.To, nil, inputStreams[0].(string))
		if err != nil {
			return err
		}
		fmt.Println("Connector successfully created!")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(connectCmd)

	// Subcommands
	connectCmd.Flags().String("to", "", "target resource name")
	connectCmd.MarkFlagRequired("to")
	connectCmd.Flags().String("from", "", "source resource name")
	connectCmd.MarkFlagRequired("from")
	connectCmd.Flags().StringP("config", "c", "", "connector configuration")
	connectCmd.Flags().String("input", "", "command delimeted list of input streams")
}
