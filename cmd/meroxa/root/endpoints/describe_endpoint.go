package endpoints

import (
	"context"
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

var (
	_ builder.CommandWithDocs    = (*DescribeEndpoint)(nil)
	_ builder.CommandWithArgs    = (*DescribeEndpoint)(nil)
	_ builder.CommandWithClient  = (*DescribeEndpoint)(nil)
	_ builder.CommandWithLogger  = (*DescribeEndpoint)(nil)
	_ builder.CommandWithExecute = (*DescribeEndpoint)(nil)
)

type describeEndpointClient interface {
	GetEndpoint(ctx context.Context, name string) (*meroxa.Endpoint, error)
}

type DescribeEndpoint struct {
	client describeEndpointClient
	logger log.Logger

	args struct {
		Name string
	}
}

func (d *DescribeEndpoint) Execute(ctx context.Context) error {
	endpoint, err := d.client.GetEndpoint(ctx, d.args.Name)
	if err != nil {
		return err
	}

	d.logger.Info(ctx, utils.EndpointsTable([]meroxa.Endpoint{*endpoint}))
	d.logger.JSON(ctx, endpoint)

	return nil
}

func (d *DescribeEndpoint) Client(client *meroxa.Client) {
	d.client = client
}

func (d *DescribeEndpoint) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *DescribeEndpoint) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires endpoint name")
	}

	d.args.Name = args[0]
	return nil
}

func (d *DescribeEndpoint) Usage() string {
	return "describe [NAME]"
}

func (d *DescribeEndpoint) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe endpoint",
	}
}
