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

package apps

import (
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/meroxa/cli/utils/display"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithAliases = (*Logs)(nil)
	_ builder.CommandWithDocs    = (*Logs)(nil)
	_ builder.CommandWithArgs    = (*Logs)(nil)
	_ builder.CommandWithClient  = (*Logs)(nil)
	_ builder.CommandWithLogger  = (*Logs)(nil)
	_ builder.CommandWithExecute = (*Logs)(nil)
)

type Logs struct {
	client applicationLogsClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}
}

type applicationLogsClient interface {
	GetApplication(ctx context.Context, nameOrUUID string) (*meroxa.Application, error)
	GetConnectorByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Connector, error)
	GetConnectorLogs(ctx context.Context, nameOrID string) (*http.Response, error)
	GetFunction(ctx context.Context, nameOrUUID string) (*meroxa.Function, error)
	GetFunctionLogs(ctx context.Context, nameOrUUID string) (*http.Response, error)
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
}

func (*Logs) Aliases() []string {
	return []string{"log"}
}

func (l *Logs) Usage() string {
	return "logs [NAMEorUUID]"
}

func (l *Logs) Docs() builder.Docs {
	return builder.Docs{
		Short:   "View relevant logs to the state of the given Turbine Data Application",
		Example: "meroxa apps logs my-turbine-application",
		Beta:    true,
	}
}

func (l *Logs) Execute(ctx context.Context) error {
	app, err := l.client.GetApplication(ctx, l.args.NameOrUUID)
	if err != nil {
		return err
	}

	connectors := make([]*display.AppExtendedConnector, 0)
	functions := make([]*meroxa.Function, 0)

	resources := app.Resources
	for _, cc := range app.Connectors {
		connector, err := l.client.GetConnectorByNameOrID(ctx, cc.Name.String)
		if err != nil {
			return err
		}

		resp, err := l.client.GetConnectorLogs(ctx, connector.Name)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return err
		}

		connectors = append(connectors, &display.AppExtendedConnector{Connector: connector, Logs: buf.String()})
	}
	for _, ff := range app.Functions {
		function, err := l.client.GetFunction(ctx, ff.Name.String)
		if err != nil {
			return err
		}

		resp, err := l.client.GetFunctionLogs(ctx, ff.Name.String)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return err
		}

		function.Logs = buf.String()
		functions = append(functions, function)
	}
	output := display.AppLogsTable(resources, connectors, functions)

	l.logger.Info(ctx, output)
	l.logger.JSON(ctx, app)

	return nil
}

func (l *Logs) Client(client meroxa.Client) {
	l.client = client
}

func (l *Logs) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *Logs) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires app name or UUID")
	}

	l.args.NameOrUUID = args[0]
	return nil
}
