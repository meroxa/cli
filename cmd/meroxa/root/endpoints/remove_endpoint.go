package endpoints

import (
	"context"
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
)

type removeEndpointClient interface {
	DeleteEndpoint(ctx context.Context, name string) error
}

type RemoveEndpoint struct {
	client removeEndpointClient
	logger log.Logger

	args struct {
		Name string
	}
}

func (r *RemoveEndpoint) Confirm(_ context.Context) (wantInput string) {
	return r.args.Name
}

func (r *RemoveEndpoint) Execute(ctx context.Context) error {
	r.logger.Infof(ctx, "Removing endpoint %q...", r.args.Name)

	err := r.client.DeleteEndpoint(ctx, r.args.Name)

	if err != nil {
		return err
	}

	r.logger.Infof(ctx, "Endpoint %q successfully removed", r.args.Name)
	// TODO: Update DeleteEndpoint in meroxa-go to return endpoint removed if needed

	return nil
}

func (r *RemoveEndpoint) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *RemoveEndpoint) Client(client *meroxa.Client) {
	r.client = client
}

func (r *RemoveEndpoint) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires endpoint name")
	}

	r.args.Name = args[0]
	return nil
}

func (r *RemoveEndpoint) Aliases() []string {
	return []string{"rm", "delete"}
}

func (r *RemoveEndpoint) Usage() string {
	return "remove NAME"
}

func (r *RemoveEndpoint) Docs() builder.Docs {
	return builder.Docs{
		Short: "Remove endpoint",
	}
}

var (
	_ builder.CommandWithDocs    = (*RemoveEndpoint)(nil)
	_ builder.CommandWithAliases = (*RemoveEndpoint)(nil)
	_ builder.CommandWithArgs    = (*RemoveEndpoint)(nil)
	_ builder.CommandWithClient  = (*RemoveEndpoint)(nil)
	_ builder.CommandWithLogger  = (*RemoveEndpoint)(nil)
	_ builder.CommandWithExecute = (*RemoveEndpoint)(nil)
	_ builder.CommandWithConfirm = (*RemoveEndpoint)(nil)
)
