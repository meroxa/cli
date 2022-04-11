/*
Copyright © 2022 Meroxa Inc

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

	"github.com/meroxa/meroxa-go/pkg/meroxa"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

var (
	_ builder.CommandWithDocs       = (*Connect)(nil)
	_ builder.CommandWithFlags      = (*Connect)(nil)
	_ builder.CommandWithLogger     = (*Connect)(nil)
	_ builder.CommandWithExecute    = (*Connect)(nil)
	_ builder.CommandWithClient     = (*Connect)(nil)
	_ builder.CommandWithDeprecated = (*Connect)(nil)
)

type Connect struct {
	logger log.Logger
	client createConnectorClient
	flags  struct {
		Source      string `long:"from" usage:"source resource name" required:"true"`
		Destination string `long:"to" usage:"destination resource name" required:"true"`
		Config      string `long:"config" usage:"connector configuration" short:"c"`
		Input       string `long:"input" usage:"command delimited list of input streams"`
		Pipeline    string `long:"pipeline" usage:"pipeline name to attach connectors to" required:"true"`
	}
}

func (c *Connect) Client(client meroxa.Client) {
	c.client = client
}

func (c *Connect) Usage() string {
	return "connect --from RESOURCE-NAME --to RESOURCE-NAME"
}

func (c *Connect) Docs() builder.Docs {
	return builder.Docs{
		Short: "Connect two resources together",
		Long: `Use the connect command to automatically configure the connectors required to pull data 
from one resource (source) to another (destination).

This command is equivalent to creating two connectors separately, 
one from the source to Meroxa and another from Meroxa to the destination:

meroxa connect --from RESOURCE-NAME --to RESOURCE-NAME --input SOURCE-INPUT --pipeline my-pipeline

or

meroxa connector create --from postgres --input accounts --pipeline my-pipeline # Creates source connector
meroxa connector create --to redshift --input orders --pipeline my-pipeline # Creates destination connector
`,
	}
}

func (c *Connect) Execute(ctx context.Context) error {
	cc := &Create{
		client: c.client,
		logger: c.logger,
	}

	cc.flags.Input = c.flags.Input
	cc.flags.Config = c.flags.Config
	cc.flags.Source = c.flags.Source
	cc.flags.Pipeline = c.flags.Pipeline

	// creates the source connector
	srcCon, err := cc.CreateConnector(ctx)

	if err != nil {
		return err
	}

	// we use the stream of the source as the input for the destination below
	inputStreams := srcCon.Streams["output"].([]interface{})
	cc.flags.Input = inputStreams[0].(string)
	cc.flags.Source = "" // unset the source to make sure cc.Create shows the right output
	cc.flags.Destination = c.flags.Destination

	destCon, err := cc.CreateConnector(ctx)

	if err != nil {
		return err
	}

	c.logger.Infof(ctx, "Source connector %q and destination connector %q successfully created!\n", srcCon.Name, destCon.Name)

	// Combine both source and destination connectors so they're included in JSON format
	connectors := []*meroxa.Connector{srcCon, destCon}

	c.logger.JSON(ctx, connectors)

	return nil
}

func (c *Connect) Flags() []builder.Flag {
	return builder.BuildFlags(&c.flags)
}

func (c *Connect) Logger(logger log.Logger) {
	c.logger = logger
}

func (*Connect) Deprecated() string {
	return "we encourage you to operate with your applications via `apps` instead."
}
