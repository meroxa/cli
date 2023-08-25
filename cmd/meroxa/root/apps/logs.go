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
	"context"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithAliases = (*Logs)(nil)
	_ builder.CommandWithDocs    = (*Logs)(nil)
	_ builder.CommandWithArgs    = (*Logs)(nil)
	_ builder.CommandWithFlags   = (*Logs)(nil)
	_ builder.CommandWithClient  = (*Logs)(nil)
	_ builder.CommandWithLogger  = (*Logs)(nil)
	_ builder.CommandWithExecute = (*Logs)(nil)
)

type Logs struct {
	client     applicationLogsClient
	logger     log.Logger
	turbineCLI turbine.CLI
	path       string

	args struct {
		NameOrUUID string
	}
	flags struct {
		Path string `long:"path" usage:"Path to the app directory (default is local directory)"`
	}
}

type applicationLogsClient interface {
	GetApplicationLogsV2(ctx context.Context, nameOrUUID string) (*meroxa.Logs, error)
	AddHeader(key, value string)
}

func (*Logs) Aliases() []string {
	return []string{"log"}
}

func (l *Logs) Usage() string {
	return `logs [NameOrUUID] [--path pwd]`
}

func (l *Logs) Flags() []builder.Flag {
	return builder.BuildFlags(&l.flags)
}

func (l *Logs) Docs() builder.Docs {
	return builder.Docs{
		Short: "View relevant logs to the state of the given Turbine Data Application",
		Long: `This command will fetch relevant logs about the Application specified in '--path'
(or current working directory if not specified) on our Meroxa Platform,
or the Application specified by the given name or UUID identifier.`,
		Example: `meroxa apps logs # assumes that the Application is in the current directory
meroxa apps logs --path /my/app
meroxa apps logs my-turbine-application`,
	}
}

func (l *Logs) Execute(ctx context.Context) error {
	var turbineLibVersion string
	nameOrUUID := l.args.NameOrUUID
	if nameOrUUID != "" && l.flags.Path != "" {
		return fmt.Errorf("supply either NameOrUUID argument or --path flag")
	}

	if nameOrUUID == "" {
		var err error
		if l.path, err = turbine.GetPath(l.flags.Path); err != nil {
			return err
		}

		config, err := turbine.ReadConfigFile(l.path)
		if err != nil {
			return err
		}
		nameOrUUID = config.Name

		if l.turbineCLI == nil {
			l.turbineCLI, err = getTurbineCLIFromLanguage(l.logger, config.Language, l.path)
			if err != nil {
				return err
			}
		}

		if turbineLibVersion, err = l.turbineCLI.GetVersion(ctx); err != nil {
			return err
		}
		addTurbineHeaders(l.client, config.Language, turbineLibVersion)
	}

	appLogs, getErr := l.client.GetApplicationLogsV2(ctx, nameOrUUID)
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
	if len(args) > 0 {
		l.args.NameOrUUID = args[0]
	}

	return nil
}
