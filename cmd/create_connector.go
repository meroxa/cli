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

type CreateConnectorClient interface {
	GetResourceByName(ctx context.Context, name string) (*meroxa.Resource, error)
	CreateConnector(ctx context.Context, input meroxa.CreateConnectorInput) (*meroxa.Connector, error)
}

type CreateConnector struct {
	input, config, metadata, source, destination, name string
	pipelineID                                         int
}

func (cc *CreateConnector) setArgs(args []string) error {
	if cc.source == "" && cc.destination == "" {
		return errors.New("requires either a source (--from) or a destination (--to)\n\nUsage:\n  meroxa create connector <custom-connector-name> [--from | --to]")
	}

	// If user specified an optional connector name
	if len(args) > 0 {
		cc.name = args[0]
	}

	return nil
}

func (cc *CreateConnector) setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&cc.input, "input", "", "", "command delimited list of input streams")
	cmd.MarkFlagRequired("input")

	cmd.Flags().StringVarP(&cc.config, "config", "c", "", "connector configuration")
	cmd.Flags().StringVarP(&cc.metadata, "metadata", "m", "", "connector metadata")
	cmd.Flags().StringVarP(&cc.source, "from", "", "", "resource name to use as source")
	cmd.Flags().StringVarP(&cc.destination, "to", "", "", "resource name to use as destination")
	cmd.Flags().IntVarP(&cc.pipelineID, "pipeline", "", 0, "ID of pipeline to attach connector to") // TODO accept name once we display it in `list connectors` command

	// Hide metadata flag for now. This could probably go away
	cmd.Flags().MarkHidden("metadata")
}

func (cc *CreateConnector) parseJSONMap(str string) (out map[string]string, err error) {
	out = make(map[string]string)
	if str != "" {
		err = json.Unmarshal([]byte(str), &out)
	}
	return out, err
}

func (cc *CreateConnector) execute(ctx context.Context, c CreateConnectorClient) (*meroxa.Connector, error) {
	config, err := cc.parseJSONMap(cc.config)
	if err != nil {
		return nil, errors.New("can't parse config, make sure it is a valid JSON map")
	}

	metadata, err := cc.parseJSONMap(cc.metadata)
	if err != nil {
		return nil, errors.New("can't parse metadata, make sure it is a valid JSON map")
	}

	// merge in input
	config["input"] = cc.input

	// merge in connector type
	var resourceName string
	switch {
	case cc.source != "":
		resourceName = cc.source
		metadata["mx:connectorType"] = "source"
	case cc.destination != "":
		resourceName = cc.destination
		metadata["mx:connectorType"] = "destination"
	default:
		return nil, errors.New("requires either a source (--from) or a destination (--to)")
	}

	res, err := c.GetResourceByName(ctx, resourceName)
	if err != nil {
		return nil, fmt.Errorf("can't fetch resource with name %q: %w", resourceName, err)
	}

	if !flagRootOutputJSON {
		switch {
		case cc.source != "":
			fmt.Printf("Creating connector from source %s...\n", resourceName)
		case cc.destination != "":
			fmt.Printf("Creating connector to destination %s...\n", resourceName)
		}
	}

	return c.CreateConnector(ctx, meroxa.CreateConnectorInput{
		Name:          cc.name,
		ResourceID:    res.ID,
		PipelineID:    cc.pipelineID,
		Configuration: config,
		Metadata:      metadata,
	})
}

func (cc *CreateConnector) output(con *meroxa.Connector) {
	if flagRootOutputJSON {
		utils.JSONPrint(con)
	} else {
		fmt.Printf("Connector %s successfully created!\n", con.Name)
	}
}

// command returns the cobra Command for `create connector`
func (cc *CreateConnector) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connector [<custom-connector-name>] [flags]",
		Short: "Create a connector",
		Long:  "Use create connector to create a connector from a source (--from) or to a destination (--to)",
		Example: "\n" +
			"meroxa create connector [<custom-connector-name>] --from pg2kafka --input accounts \n" +
			"meroxa create connector [<custom-connector-name>] --to pg2redshift --input orders # --input will be the desired stream",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cc.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()

			if err != nil {
				return err
			}

			res, err := cc.execute(ctx, c)

			if err != nil {
				return err
			}

			cc.output(res)

			return nil
		},
	}

	cc.setFlags(cmd)

	return cmd
}
