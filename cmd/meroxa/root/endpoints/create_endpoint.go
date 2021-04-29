package endpoints

import (
	"context"

	"github.com/meroxa/meroxa-go"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type createEndpointClient interface {
	CreateEndpoint(ctx context.Context, name, protocol, stream string) error
}

type CreateEndpoint struct {
	client createEndpointClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Protocol string `long:"protocol" short:"p" usage:"protocol, value can be http or grpc" required:"true"`
		Stream   string `long:"stream" short:"s" usage:"stream name" required:"true"`
	}
}

func (c *CreateEndpoint) Execute(ctx context.Context) error {
	c.logger.Info(ctx, "Creating endpoint...")

	err := c.client.CreateEndpoint(ctx, c.args.Name, c.flags.Protocol, c.flags.Stream)

	if err != nil {
		return err
	}

	c.logger.Info(ctx, "Endpoint successfully created!")

	return nil
}

func (c *CreateEndpoint) Client(client *meroxa.Client) {
	c.client = client
}

func (c *CreateEndpoint) Logger(logger log.Logger) {
	c.logger = logger
}

func (c *CreateEndpoint) Flags() []builder.Flag {
	return builder.BuildFlags(&c.flags)
}

func (c *CreateEndpoint) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Create an endpoint",
		Long:    "Use create endpoint to expose an endpoint to a connector stream",
		Example: "meroxa endpoints create my-endpoint --protocol http --stream my-stream",
	}
}

func (c *CreateEndpoint) Usage() string {
	return "create [NAME] [flags]"
}

func (c *CreateEndpoint) ParseArgs(args []string) error {
	if len(args) > 0 {
		c.args.Name = args[0]
	}
	return nil
}

var (
	_ builder.CommandWithDocs    = (*CreateEndpoint)(nil)
	_ builder.CommandWithArgs    = (*CreateEndpoint)(nil)
	_ builder.CommandWithFlags   = (*CreateEndpoint)(nil)
	_ builder.CommandWithClient  = (*CreateEndpoint)(nil)
	_ builder.CommandWithLogger  = (*CreateEndpoint)(nil)
	_ builder.CommandWithExecute = (*CreateEndpoint)(nil)
)
