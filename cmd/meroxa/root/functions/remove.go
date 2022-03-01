package functions

import (
	"context"
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs             = (*Remove)(nil)
	_ builder.CommandWithAliases          = (*Remove)(nil)
	_ builder.CommandWithArgs             = (*Remove)(nil)
	_ builder.CommandWithClient           = (*Remove)(nil)
	_ builder.CommandWithLogger           = (*Remove)(nil)
	_ builder.CommandWithExecute          = (*Remove)(nil)
	_ builder.CommandWithConfirmWithValue = (*Remove)(nil)
)

type removeFunctionClient interface {
	DeleteFunction(ctx context.Context, appNameOrUUID, nameOrUUID string) (*meroxa.Function, error)
}

type Remove struct {
	client removeFunctionClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}

	flags struct {
		Application string `long:"app" usage:"application name or UUID to which this function belongs" required:"true"`
	}
}

func (r *Remove) Usage() string {
	return "remove NAMEorUUID"
}

func (r *Remove) Docs() builder.Docs {
	return builder.Docs{
		Short: "Remove function",
	}
}

func (r *Remove) ValueToConfirm(_ context.Context) (wantInput string) {
	return r.args.NameOrUUID
}

func (r *Remove) Execute(ctx context.Context) error {
	r.logger.Infof(ctx, "Function %q is being removed...", r.args.NameOrUUID)

	e, err := r.client.DeleteFunction(ctx, r.flags.Application, r.args.NameOrUUID)
	if err != nil {
		return err
	}

	r.logger.Infof(ctx, "Function %q successfully removed", r.args.NameOrUUID)
	r.logger.JSON(ctx, e)

	return nil
}

func (r *Remove) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Remove) Client(client meroxa.Client) {
	r.client = client
}

func (r *Remove) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires function name")
	}

	r.args.NameOrUUID = args[0]
	return nil
}

func (r *Remove) Aliases() []string {
	return []string{"rm", "delete"}
}
