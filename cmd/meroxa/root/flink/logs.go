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

package flink

import (
	"context"
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
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
	client flinkLogsClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}
}

type flinkLogsClient interface {
	GetFlinkLogsV2(ctx context.Context, nameOrUUID string) (*meroxa.Logs, error)
}

func (*Logs) Aliases() []string {
	return []string{"log"}
}

func (l *Logs) Usage() string {
	return `logs [NameOrUUID] [--path pwd]`
}

func (l *Logs) Docs() builder.Docs {
	return builder.Docs{
		Short: "View relevant logs to the state of the given Flink Job",
		Example: `meroxa jobs logs my-flink-job-name
meroxa jobs logs my-flink-job-uuid`,
	}
}

func (l *Logs) Execute(ctx context.Context) error {
	nameOrUUID := l.args.NameOrUUID

	appLogs, getErr := l.client.GetFlinkLogsV2(ctx, nameOrUUID)
	if getErr != nil {
		return getErr
	}

	output := display.LogsTable(appLogs)

	l.logger.Info(ctx, output)
	l.logger.JSON(ctx, appLogs)

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
		return errors.New("requires Flink Job name or UUID")
	}

	l.args.NameOrUUID = args[0]
	return nil
}
