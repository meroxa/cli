package functions

import (
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Logs)(nil)
	_ builder.CommandWithArgs    = (*Logs)(nil)
	_ builder.CommandWithClient  = (*Logs)(nil)
	_ builder.CommandWithLogger  = (*Logs)(nil)
	_ builder.CommandWithExecute = (*Logs)(nil)
)

type functionLogsClient interface {
	GetFunctionLogs(ctx context.Context, nameOrID string) (*http.Response, error)
}

type Logs struct {
	client functionLogsClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}
}

func (l *Logs) Usage() string {
	return "logs NAME"
}

func (l *Logs) Docs() builder.Docs {
	return builder.Docs{
		Short: "Print logs for a function",
	}
}

func (l *Logs) Execute(ctx context.Context) error {
	resp, err := l.client.GetFunctionLogs(ctx, l.args.NameOrID)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	if err != nil {
		return err
	}

	l.logger.Info(ctx, buf.String())

	return nil
}

func (l *Logs) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *Logs) Client(client meroxa.Client) {
	l.client = client
}

func (l *Logs) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires function name")
	}

	l.args.NameOrID = args[0]
	return nil
}
