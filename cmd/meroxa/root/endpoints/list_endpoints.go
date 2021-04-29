package endpoints

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

var (
	_ builder.CommandWithDocs    = (*ListEndpoints)(nil)
	_ builder.CommandWithClient  = (*ListEndpoints)(nil)
	_ builder.CommandWithLogger  = (*ListEndpoints)(nil)
	_ builder.CommandWithExecute = (*ListEndpoints)(nil)
	_ builder.CommandWithAliases = (*ListEndpoints)(nil)
)

type listEndpointsClient interface {
	ListEndpoints(ctx context.Context) ([]meroxa.Endpoint, error)
}

type ListEndpoints struct {
	client listEndpointsClient
	logger log.Logger
}

func (l *ListEndpoints) Aliases() []string {
	return []string{"ls"}
}

func (l *ListEndpoints) Execute(ctx context.Context) error {
	var err error
	endpoints, err := l.client.ListEndpoints(ctx)
	if err != nil {
		return err
	}

	l.logger.JSON(ctx, endpoints)
	l.logger.Info(ctx, utils.EndpointsTable(endpoints))

	return nil
}

func (l *ListEndpoints) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *ListEndpoints) Client(client *meroxa.Client) {
	l.client = client
}

func (l *ListEndpoints) Usage() string {
	return "list"
}

func (l *ListEndpoints) Docs() builder.Docs {
	return builder.Docs{
		Short: "List endpoints",
	}
}
