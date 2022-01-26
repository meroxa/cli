package functions

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs      = (*List)(nil)
	_ builder.CommandWithClient    = (*List)(nil)
	_ builder.CommandWithLogger    = (*List)(nil)
	_ builder.CommandWithExecute   = (*List)(nil)
	_ builder.CommandWithAliases   = (*List)(nil)
	_ builder.CommandWithNoHeaders = (*List)(nil)
)

type listFunctionClient interface {
	ListFunctions(ctx context.Context) ([]*meroxa.Function, error)
}

type List struct {
	client      listFunctionClient
	logger      log.Logger
	hideHeaders bool
}

func (l *List) Execute(ctx context.Context) error {
	funs, err := l.client.ListFunctions(ctx)
	if err != nil {
		return err
	}

	l.logger.JSON(ctx, funs)
	l.logger.Info(ctx, utils.FunctionsTable(funs, l.hideHeaders))

	return nil
}

func (l *List) Usage() string {
	return "list"
}

func (l *List) Docs() builder.Docs {
	return builder.Docs{
		Short: "List functions",
	}
}

func (l *List) Aliases() []string {
	return []string{"ls"}
}

func (l *List) Client(client meroxa.Client) {
	l.client = client
}

func (l *List) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *List) HideHeaders(hide bool) {
	l.hideHeaders = hide
}
