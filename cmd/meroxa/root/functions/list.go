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
	ListApplications(ctx context.Context) ([]*meroxa.Application, error)
	ListFunctions(ctx context.Context, appNameOrUUID string) ([]*meroxa.Function, error)
}

type List struct {
	client      listFunctionClient
	logger      log.Logger
	hideHeaders bool

	flags struct {
		Application string `long:"app" usage:"application name or UUID to which this function belongs"`
	}
}

func (l *List) Execute(ctx context.Context) error {
	var err error
	funs := make([]*meroxa.Function, 0)
	if l.flags.Application != "" {
		funs, err = l.client.ListFunctions(ctx, l.flags.Application)
		if err != nil {
			return err
		}
	} else {
		apps, err := l.client.ListApplications(ctx)
		if err != nil {
			return err
		}
		for _, app := range apps {
			fs, err := l.client.ListFunctions(ctx, app.UUID)
			if err != nil {
				return err
			}
			funs = append(funs, fs...)
		}
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
