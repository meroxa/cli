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
	_ builder.CommandWithDocs    = (*Describe)(nil)
	_ builder.CommandWithArgs    = (*Describe)(nil)
	_ builder.CommandWithFlags   = (*Describe)(nil)
	_ builder.CommandWithClient  = (*Describe)(nil)
	_ builder.CommandWithLogger  = (*Describe)(nil)
	_ builder.CommandWithExecute = (*Describe)(nil)
)

type describeApplicationClient interface {
	GetApplication(ctx context.Context, nameOrUUID string) (*meroxa.Application, error)
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
	GetConnectorByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Connector, error)
	GetFunction(ctx context.Context, nameOrUUID string) (*meroxa.Function, error)
	AddHeader(key, value string)
}

type Describe struct {
	client     describeApplicationClient
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

func (d *Describe) Usage() string {
	return "describe [NameOrUUID] [--path pwd]"
}

func (d *Describe) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe a Turbine Data Application",
		Long: `This command will fetch details about the Application specified in '--path'
(or current working directory if not specified) on our Meroxa Platform,
or the Application specified by the given name or UUID identifier.`,
		Example: `meroxa apps describe # assumes that the Application is in the current directory
meroxa apps describe --path /my/app
meroxa apps describe NAMEorUUID`,
	}
}

func (d *Describe) Execute(ctx context.Context) error {
	var turbineLibVersion string
	nameOrUUID := d.args.NameOrUUID
	if nameOrUUID != "" && d.flags.Path != "" {
		return fmt.Errorf("supply either NameOrUUID argument or --path flag")
	}

	if nameOrUUID == "" {
		var err error
		if d.path, err = turbine.GetPath(d.flags.Path); err != nil {
			return err
		}

		config, err := turbine.ReadConfigFile(d.path)
		if err != nil {
			return err
		}
		nameOrUUID = config.Name

		if d.turbineCLI == nil {
			d.turbineCLI, err = getTurbineCLIFromLanguage(d.logger, config.Language, d.path)
			if err != nil {
				return err
			}
		}

		if turbineLibVersion, err = d.turbineCLI.GetVersion(ctx); err != nil {
			return err
		}
		addTurbineHeaders(d.client, config.Language, turbineLibVersion)
	}

	app, err := d.client.GetApplication(ctx, nameOrUUID)
	if err != nil {
		return err
	}

	d.logger.Info(ctx, display.AppTable(app))
	d.logger.JSON(ctx, app)

	dashboardURL := fmt.Sprintf("https://dashboard.meroxa.io/apps/%s/detail", app.Name)
	d.logger.Info(ctx, fmt.Sprintf("\n ✨ To visualize your application, visit %s", dashboardURL))
	return nil
}

func (d *Describe) Client(client meroxa.Client) {
	d.client = client
}

func (d *Describe) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Describe) ParseArgs(args []string) error {
	if len(args) > 0 {
		d.args.NameOrUUID = args[0]
	}

	return nil
}
