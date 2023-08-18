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

package builds

import (
	"context"
	"errors"
	"net/http"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Logs)(nil)
	_ builder.CommandWithAliases = (*Logs)(nil)
	_ builder.CommandWithArgs    = (*Logs)(nil)
	_ builder.CommandWithClient  = (*Logs)(nil)
	_ builder.CommandWithLogger  = (*Logs)(nil)
	_ builder.CommandWithExecute = (*Logs)(nil)
)

type buildLogsClient interface {
	GetBuildLogs(ctx context.Context, uuid string) (*http.Response, error)
	GetBuildLogsV2(ctx context.Context, uuid string) (*meroxa.Logs, error)
}

type Logs struct {
	client buildLogsClient
	logger log.Logger

	args struct {
		UUID string
	}
}

func (l *Logs) Usage() string {
	return "logs [UUID]"
}

func (*Logs) Aliases() []string {
	return []string{"log"}
}

func (l *Logs) Docs() builder.Docs {
	return builder.Docs{
		Short: "List a Meroxa Process Build's Logs",
	}
}

func (l *Logs) Execute(ctx context.Context) error {
	buildLogs, getErr := l.client.GetBuildLogsV2(ctx, l.args.UUID)
	if getErr != nil {
		return getErr
	}

	output := display.BuildsLogsTable(buildLogs)

	l.logger.Info(ctx, output)
	l.logger.JSON(ctx, buildLogs)

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
		return errors.New("requires build UUID")
	}

	l.args.UUID = args[0]
	return nil
}
