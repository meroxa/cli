package connectors

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/meroxa/meroxa-go"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

var (
	_ builder.CommandWithDocs    = (*LogsConnector)(nil)
	_ builder.CommandWithArgs    = (*LogsConnector)(nil)
	_ builder.CommandWithClient  = (*LogsConnector)(nil)
	_ builder.CommandWithLogger  = (*LogsConnector)(nil)
	_ builder.CommandWithExecute = (*LogsConnector)(nil)
)

type logsConnectorClient interface {
	GetConnectorLogs(ctx context.Context, connectorName string) (*http.Response, error)
}

type LogsConnector struct {
	client logsConnectorClient
	logger log.Logger

	args struct {
		Name string
	}
}

func (l *LogsConnector) Execute(ctx context.Context) error {
	resp, err := l.client.GetConnectorLogs(ctx, l.args.Name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err
	}

	os.Stdout.Write([]byte("\n"))

	return nil
}

func (l *LogsConnector) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *LogsConnector) Client(client *meroxa.Client) {
	l.client = client
}

func (l *LogsConnector) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires connector name")
	}

	l.args.Name = args[0]
	return nil
}

func (l *LogsConnector) Usage() string {
	return "logs NAME"
}

func (l *LogsConnector) Docs() builder.Docs {
	return builder.Docs{
		Short: "Print logs for a connector",
	}
}
