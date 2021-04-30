package connectors

import (
	"context"
	"errors"

	"github.com/meroxa/cli/utils"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
)

var (
	_ builder.CommandWithDocs    = (*DescribeConnector)(nil)
	_ builder.CommandWithArgs    = (*DescribeConnector)(nil)
	_ builder.CommandWithClient  = (*DescribeConnector)(nil)
	_ builder.CommandWithLogger  = (*DescribeConnector)(nil)
	_ builder.CommandWithExecute = (*DescribeConnector)(nil)
)

type describeConnectorClient interface {
	GetConnectorByName(ctx context.Context, name string) (*meroxa.Connector, error)
}

type DescribeConnector struct {
	client describeConnectorClient
	logger log.Logger

	args struct {
		Name string
	}
}

func (d *DescribeConnector) Usage() string {
	return "describe [NAME]"
}

func (d *DescribeConnector) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe connector",
	}
}

func (d *DescribeConnector) Execute(ctx context.Context) error {
	connector, err := d.client.GetConnectorByName(ctx, d.args.Name)
	if err != nil {
		return err
	}

	d.logger.Info(ctx, utils.ConnectorsTable([]*meroxa.Connector{connector}))
	d.logger.JSON(ctx, connector)

	return nil
}

func (d *DescribeConnector) Client(client *meroxa.Client) {
	d.client = client
}

func (d *DescribeConnector) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *DescribeConnector) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires connector name")
	}

	d.args.Name = args[0]
	return nil
}
