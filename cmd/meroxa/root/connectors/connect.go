package connectors

import (
	"context"
	"os"

	"github.com/meroxa/meroxa-go"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Connect struct {
	logger log.Logger
	client createConnectorClient
	flags  struct {
		Source      string `long:"from" usage:"source resource name" required:"true"`
		Destination string `long:"to" usage:"destination resource name" required:"true"`
		Config      string `long:"config" usage:"connector configuration" short:"c"`
		Input       string `long:"input" usage:"command delimited list of input streams"`
		Pipeline    string `long:"pipeline" short:"" usage:"pipeline name to attach connectors to"`
	}
}

func (c *Connect) Client(client *meroxa.Client) {
	c.client = client
}

func (c *Connect) Usage() string {
	return "connect --from RESOURCE-NAME --to RESOURCE-NAME"
}

func (c *Connect) Docs() builder.Docs {
	longOutput := `Use the connect command to automatically configure the connectors required to pull data 
from one resource (source) to another (destination).

This command is equivalent to creating two connectors separately, 
one from the source to Meroxa and another from Meroxa to the destination:

meroxa connect --from RESOURCE-NAME --to RESOURCE-NAME --input SOURCE-INPUT

or
`
	// Adapt the output based on the CLI version
	if _, ok := os.LookupEnv("MEROXA_V2"); ok {
		longOutput += `
meroxa connector create --from postgres --input accounts # Creates source connector
meroxa connector create --to redshift --input orders # Creates destination connector
`
	} else {
		longOutput += `
meroxa create connector --from postgres --input accounts # Creates source connector
meroxa create connector --to redshift --input orders # Creates destination connector
`
	}

	return builder.Docs{
		Short: "Connect two resources together",
		Long:  longOutput,
	}
}

func (c *Connect) Execute(ctx context.Context) error {
	cc := &CreateConnector{
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
	cc.flags.Source = "" // unset the source to make sure cc.CreateConnector shows the right output
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

var (
	_ builder.CommandWithDocs    = (*Connect)(nil)
	_ builder.CommandWithFlags   = (*Connect)(nil)
	_ builder.CommandWithLogger  = (*Connect)(nil)
	_ builder.CommandWithExecute = (*Connect)(nil)
	_ builder.CommandWithClient  = (*Connect)(nil)
)
