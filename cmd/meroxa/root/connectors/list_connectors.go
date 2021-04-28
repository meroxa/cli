package connectors

import (
	"context"

	"github.com/meroxa/cli/utils"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
)

var (
	_ builder.CommandWithDocs    = (*ListConnectors)(nil)
	_ builder.CommandWithClient  = (*ListConnectors)(nil)
	_ builder.CommandWithLogger  = (*ListConnectors)(nil)
	_ builder.CommandWithExecute = (*ListConnectors)(nil)
	_ builder.CommandWithFlags   = (*ListConnectors)(nil)
)

type listConnectorsClient interface {
	ListConnectors(ctx context.Context) ([]*meroxa.Connector, error)
	ListPipelineConnectors(ctx context.Context, pipelineID int) ([]*meroxa.Connector, error)
	GetPipelineByName(ctx context.Context, name string) (*meroxa.Pipeline, error)
}

type ListConnectors struct {
	client listConnectorsClient
	logger log.Logger

	flags struct {
		Pipeline string `long:"pipeline"        short:""  usage:"filter connectors by pipeline name"        required:"false"`
	}
}

func (l *ListConnectors) Execute(ctx context.Context) error {
	var err error
	var connectors []*meroxa.Connector

	// Filtering by pipeline name
	if l.flags.Pipeline != "" {
		var p *meroxa.Pipeline

		p, err = l.client.GetPipelineByName(ctx, l.flags.Pipeline)

		if err != nil {
			return err
		}

		connectors, err = l.client.ListPipelineConnectors(ctx, p.ID)

		if err != nil {
			return err
		}
	} else {
		connectors, err = l.client.ListConnectors(ctx)

		if err != nil {
			return err
		}
	}

	l.logger.JSON(ctx, connectors)
	l.logger.Info(ctx, utils.ConnectorsTable(connectors))

	return nil
}

func (l *ListConnectors) Flags() []builder.Flag {
	return builder.BuildFlags(&l.flags)
}

func (l *ListConnectors) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *ListConnectors) Client(client *meroxa.Client) {
	l.client = client
}

func (l *ListConnectors) Usage() string {
	return "list"
}

func (l *ListConnectors) Docs() builder.Docs {
	return builder.Docs{
		Short: "List connectors",
	}
}
