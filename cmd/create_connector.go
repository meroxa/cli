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

// Config defines a dictionary to be used across our cmd package
type Config map[string]string

// Set a key value pair
func (c Config) Set(key, value string) error {
	c[key] = value
	return nil
}

// Get a key value pair
func (c Config) Get(key string) (string, bool) {
	v, ok := c[key]
	return v, ok
}

// Merge one Config definition onto another
func (c Config) Merge(cfg Config) error {
	for k, v := range cfg {
		_, exist := c[k]
		if exist {
			return fmt.Errorf("merge config, key %s already present", k)
		}
		c[k] = v
	}
	return nil
}

// CreateConnectorCmd represents the `meroxa create connector` command
func CreateConnectorCmd() *cobra.Command {
	createConnectorCmd := &cobra.Command{
		Use:   "connector [<custom-connector-name>] [flags]",
		Short: "Create a connector",
		Long:  "Use create connector to create a connector from a source (--from) or to a destination (--to)",
		Example: "\n" +
			"meroxa create connector [<custom-connector-name>] --from pg2kafka --input accounts \n" +
			"meroxa create connector [<custom-connector-name>] --to pg2redshift --input orders # --input will be the desired stream",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if source == "" && destination == "" {
				return errors.New("requires either a source (--from) or a destination (--to)\n\nUsage:\n  meroxa create connector <custom-connector-name> [--from | --to]")
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &Config{}
			if cfgString != "" {
				err := json.Unmarshal([]byte(cfgString), cfg)
				if err != nil {
					return err
				}
			}

			// Process metadata
			metadata := map[string]string{}
			if metadataString != "" {
				err := json.Unmarshal([]byte(metadataString), &metadata)
				if err != nil {
					return err
				}
			}

			// merge in input
			err := cfg.Set("input", input)
			if err != nil {
				return err
			}

			if source != "" {
				res = source
				metadata["mx:connectorType"] = "source"
			} else if destination != "" {
				res = destination
				metadata["mx:connectorType"] = "destination"
			}

			// If user specified an optional connector name
			if len(args) > 0 {
				con = args[0]
			}

			if !flagRootOutputJSON {
				fmt.Println("Creating connector...")
			}

			c, err := createConnector(con, res, cfg, metadata, input)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				utils.JSONPrint(c)
			} else {
				fmt.Printf("Connector %s successfully created!\n", c.Name)
			}

			return nil
		},
	}

	createConnectorCmd.Flags().StringVarP(&cfgString, "config", "c", "", "connector configuration")
	createConnectorCmd.Flags().StringVarP(&metadataString, "metadata", "m", "", "connector metadata")
	createConnectorCmd.Flags().StringVarP(&input, "input", "", "", "command delimeted list of input streams")
	createConnectorCmd.MarkFlagRequired("input")
	createConnectorCmd.Flags().StringVarP(&source, "from", "", "", "resource name to use as source")
	createConnectorCmd.Flags().StringVarP(&destination, "to", "", "", "resource name to use as destination")

	// Hide metadata flag for now. This could probably go away
	createConnectorCmd.Flags().MarkHidden("metadata")

	return createConnectorCmd
}

func createConnector(connectorName string, resourceName string, config *Config, metadata map[string]string, input string) (*meroxa.Connector, error) {
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

	cfg := Config{}
	if config != nil {
		cfg.Merge(*config)
	}

	// merge in input
	if input != "" {
		err = cfg.Set("input", input)
		if err != nil {
			return nil, err
		}
	}

	con, err := c.CreateConnector(ctx, connectorName, res.ID, cfg, metadata)
	if err != nil {
		return nil, err
	}

	return con, nil
}
