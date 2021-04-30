package connectors

import (
	"context"
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/meroxa-go"
)

type removeConnectorClient interface {
	GetConnectorByName(ctx context.Context, name string) (*meroxa.Connector, error)
	DeleteConnector(ctx context.Context, id int) error
}

type RemoveConnector struct {
	client removeConnectorClient
	logger log.Logger

	args struct {
		Name string
	}
}

func (r *RemoveConnector) Usage() string {
	return "remove NAME"
}

func (r *RemoveConnector) Docs() builder.Docs {
	return builder.Docs{
		Short: "Remove connector",
	}
}

func (r *RemoveConnector) Confirm(_ context.Context) (wantInput string) {
	return r.args.Name
}

func (r *RemoveConnector) Execute(ctx context.Context) error {
	r.logger.Infof(ctx, "Removing connector %q...", r.args.Name)

	con, err := r.client.GetConnectorByName(ctx, r.args.Name)
	if err != nil {
		return err
	}

	err = r.client.DeleteConnector(ctx, con.ID)

	if err != nil {
		return err
	}

	r.logger.Infof(ctx, "Connector %q successfully removed", r.args.Name)
	r.logger.JSON(ctx, con)

	return nil
}

func (r *RemoveConnector) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *RemoveConnector) Client(client *meroxa.Client) {
	r.client = client
}

func (r *RemoveConnector) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires connector name")
	}

	r.args.Name = args[0]
	return nil
}

func (r *RemoveConnector) Aliases() []string {
	return []string{"rm", "delete"}
}

var (
	_ builder.CommandWithDocs    = (*RemoveConnector)(nil)
	_ builder.CommandWithAliases = (*RemoveConnector)(nil)
	_ builder.CommandWithArgs    = (*RemoveConnector)(nil)
	_ builder.CommandWithClient  = (*RemoveConnector)(nil)
	_ builder.CommandWithLogger  = (*RemoveConnector)(nil)
	_ builder.CommandWithExecute = (*RemoveConnector)(nil)
	_ builder.CommandWithConfirm = (*RemoveConnector)(nil)
)
