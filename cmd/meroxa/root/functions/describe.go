package functions

import (
	"context"
	"errors"

	"github.com/meroxa/cli/utils/display"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Describe)(nil)
	_ builder.CommandWithArgs    = (*Describe)(nil)
	_ builder.CommandWithClient  = (*Describe)(nil)
	_ builder.CommandWithLogger  = (*Describe)(nil)
	_ builder.CommandWithExecute = (*Describe)(nil)
)

type describeFunctionClient interface {
	GetFunction(ctx context.Context, nameOrUUID string) (*meroxa.Function, error)
}

type Describe struct {
	client describeFunctionClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}
}

func (d *Describe) Usage() string {
	return "describe [NAMEorUUID]"
}

func (d *Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe function",
	}
}

func (d *Describe) Execute(ctx context.Context) error {
	fun, err := d.client.GetFunction(ctx, d.args.NameOrUUID)
	if err != nil {
		return err
	}

	d.logger.Info(ctx, display.FunctionTable(fun))
	d.logger.JSON(ctx, fun)

	return nil
}

func (d *Describe) Client(client meroxa.Client) {
	d.client = client
}

func (d *Describe) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Describe) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires function name")
	}

	d.args.NameOrUUID = args[0]
	return nil
}
