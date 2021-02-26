package cmd

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

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/meroxa/cli/display"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

var (
	con            string // connector name
	res            string // resource name
	cfgString      string
	metadataString string
	input          string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Meroxa pipeline components",
	Long: `Use the create command to create various Meroxa pipeline components
including connectors.`,
}

var createConnectorCmd = &cobra.Command{
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
		// Custom connector name (which is optional)
		if len(args) > 0 {
			con = args[0]
		}

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

		if !flagRootOutputJSON {
			fmt.Println("Creating connector...")
		}

		if source != "" {
			metadata["mx:connectorType"] = "source"
		} else if destination != "" {
			metadata["mx:connectorType"] = "destination"
		}

		c, err := createConnector(con, res, cfg, metadata, input)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(c)
		} else {
			fmt.Println("Connector successfully created!")
			display.PrettyPrint("connector", c)
		}

		return nil
	},
}

var createPipelineCmd = &cobra.Command{
	Use:   "pipeline <name>",
	Short: "Create a pipeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a pipeline name\n\nUsage:\n  meroxa create pipeline <name> [flags]")
		}
		pipelineName := args[0]

		c, err := client()
		if err != nil {
			return err
		}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
			display.JSONPrint(res)
		} else {
			fmt.Println("Pipeline successfully created!")
			display.PrettyPrint("pipeline", res)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.AddCommand(createConnectorCmd)
	createConnectorCmd.Flags().StringVarP(&cfgString, "config", "c", "", "connector configuration")
	createConnectorCmd.Flags().StringVarP(&metadataString, "metadata", "m", "", "connector metadata")
	createConnectorCmd.Flags().StringVarP(&input, "input", "", "", "command delimeted list of input streams")
	createConnectorCmd.MarkFlagRequired("input")
	createConnectorCmd.Flags().StringVarP(&source, "from", "", "", "resource name to use as source")
	createConnectorCmd.Flags().StringVarP(&destination, "to", "", "", "resource name to use as destination")

	// Hide metadata flag for now. This could probably go away
	createConnectorCmd.Flags().MarkHidden("metadata")

	createCmd.AddCommand(createPipelineCmd)
	createPipelineCmd.Flags().StringP("metadata", "m", "", "pipeline metadata")
}

func createConnector(connectorName string, resourceName string, config *Config, metadata map[string]string, input string) (*meroxa.Connector, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}

	// get resource ID from name
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.GetResourceByName(ctx, resourceName)
	if err != nil {
		return nil, err
	}

	// create connector
	ctx = context.Background()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
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
