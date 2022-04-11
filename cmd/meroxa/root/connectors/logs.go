/*
Copyright © 2022 Meroxa Inc

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

	"github.com/meroxa/meroxa-go/pkg/meroxa"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

var (
	_ builder.CommandWithDocs       = (*Logs)(nil)
	_ builder.CommandWithArgs       = (*Logs)(nil)
	_ builder.CommandWithClient     = (*Logs)(nil)
	_ builder.CommandWithLogger     = (*Logs)(nil)
	_ builder.CommandWithExecute    = (*Logs)(nil)
	_ builder.CommandWithDeprecated = (*Logs)(nil)
)

type logsConnectorClient interface {
	GetConnectorLogs(ctx context.Context, nameOrID string) (*http.Response, error)
}

type Logs struct {
	client logsConnectorClient
	logger log.Logger

	args struct {
		NameOrID string
	}
}

func (l *Logs) Deprecated() string {
	return "we encourage you to operate with your applications via `meroxa apps` instead."
}

func (l *Logs) Usage() string {
	return "logs NAME"
}

func (l *Logs) Docs() builder.Docs {
	return builder.Docs{
		Short: "Print logs for a connector",
	}
}

func (l *Logs) Execute(ctx context.Context) error {
	resp, err := l.client.GetConnectorLogs(ctx, l.args.NameOrID)

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
		return errors.New("requires connector name")
	}

	l.args.NameOrID = args[0]
	return nil
}
