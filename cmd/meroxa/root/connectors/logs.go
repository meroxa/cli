/*
Copyright © 2021 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package connectors

import (
	"bytes"
	"context"
	"errors"
	"net/http"

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

func (l *LogsConnector) Usage() string {
	return "logs NAME"
}

func (l *LogsConnector) Docs() builder.Docs {
	return builder.Docs{
		Short: "Print logs for a connector",
	}
}

func (l *LogsConnector) Execute(ctx context.Context) error {
	resp, err := l.client.GetConnectorLogs(ctx, l.args.Name)

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
