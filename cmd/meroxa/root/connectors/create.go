/*
Copyright © 2021 Meroxa Inc

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

package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Create)(nil)
	_ builder.CommandWithArgs    = (*Create)(nil)
	_ builder.CommandWithFlags   = (*Create)(nil)
	_ builder.CommandWithClient  = (*Create)(nil)
	_ builder.CommandWithLogger  = (*Create)(nil)
	_ builder.CommandWithExecute = (*Create)(nil)
)

type createConnectorClient interface {
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
	CreateConnector(ctx context.Context, input *meroxa.CreateConnectorInput) (*meroxa.Connector, error)
}

type Create struct {
	client createConnectorClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Input       string `long:"input" usage:"command delimited list of input streams"`
		Config      string `long:"config" short:"c" usage:"connector configuration"`
		Metadata    string `long:"metadata" short:"m" usage:"connector metadata" hidden:"true"`
		Source      string `long:"from" usage:"resource name to use as source"`
		Destination string `long:"to" usage:"resource name to use as destination"`
		Pipeline    string `long:"pipeline" usage:"pipeline name to attach connector to" required:"true"`
	}
}

func (c *Create) Usage() string {
	return "create [NAME] [flags]"
}

func (c *Create) Docs() builder.Docs {
	return builder.Docs{
		Short: "Create a connector",
		Long:  "Use `connectors create` to create a connector from a source (--from) or to a destination (--to) within a pipeline (--pipeline)",
		Example: "\n" +
			"meroxa connectors create [NAME] --from pg2kafka --input accounts --pipeline my-pipeline\n" +
			"meroxa connectors create [NAME] --to pg2redshift --input orders --pipeline my-pipeline # --input will be the desired stream\n" +
			"meroxa connectors create [NAME] --to pg2redshift --input orders --pipeline my-pipeline\n",
	}
}

func (c *Create) parseJSONMap(str string) (out map[string]interface{}, err error) {
	out = make(map[string]interface{})
	if str != "" {
		err = json.Unmarshal([]byte(str), &out)
	}
	return out, err
}

func (c *Create) CreateConnector(ctx context.Context) (*meroxa.Connector, error) {
	config, err := c.parseJSONMap(c.flags.Config)
	if err != nil {
		return nil, errors.New("can't parse config, make sure it is a valid JSON map")
	}

	metadata, err := c.parseJSONMap(c.flags.Metadata)
	if err != nil {
		return nil, fmt.Errorf("could not parse metadata: %w", err)
	}

	if c.flags.Pipeline == "" {
		return nil, errors.New("requires pipeline name (--pipeline)")
	}

	var connectorType meroxa.ConnectorType
	// merge in connector type
	var resourceName string
	switch {
	case c.flags.Source != "":
		resourceName = c.flags.Source
		connectorType = meroxa.ConnectorTypeSource
	case c.flags.Destination != "":
		resourceName = c.flags.Destination
		connectorType = meroxa.ConnectorTypeDestination
	default:
		return nil, errors.New("requires either a source (--from) or a destination (--to)")
	}

	res, err := c.client.GetResourceByNameOrID(ctx, resourceName)
	if err != nil {
		return nil, fmt.Errorf("can't fetch resource with name %q: %w", resourceName, err)
	}

	switch {
	case c.flags.Source != "":
		c.logger.Infof(ctx, "Creating connector from source %q in pipeline %q...\n", resourceName, c.flags.Pipeline)
	case c.flags.Destination != "":
		c.logger.Infof(ctx, "Creating connector to destination %q in pipeline %q...\n", resourceName, c.flags.Pipeline)
	}

	ci := &meroxa.CreateConnectorInput{
		Name:          c.args.Name,
		ResourceID:    res.ID,
		PipelineName:  c.flags.Pipeline,
		Configuration: config,
		Metadata:      metadata,
		Type:          connectorType,
		Input:         c.flags.Input,
	}

	return c.client.CreateConnector(ctx, ci)
}

func (c *Create) Execute(ctx context.Context) error {
	// TODO: Implement something like dependent flags in Builder
	if c.flags.Source == "" && c.flags.Destination == "" {
		return errors.New("requires either a source (--from) or a destination (--to)")
	}

	connector, err := c.CreateConnector(ctx)

	if err != nil {
		return err
	}

	c.logger.Infof(ctx, "Connector %q successfully created!\n", connector.Name)
	c.logger.JSON(ctx, connector)

	return nil
}

func (c *Create) Client(client meroxa.Client) {
	c.client = client
}

func (c *Create) Logger(logger log.Logger) {
	c.logger = logger
}

func (c *Create) Flags() []builder.Flag {
	return builder.BuildFlags(&c.flags)
}

func (c *Create) ParseArgs(args []string) error {
	if len(args) > 0 {
		c.args.Name = args[0]
	}
	return nil
}
