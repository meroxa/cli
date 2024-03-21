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
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
)

var (
	_ builder.CommandWithDocs        = (*Describe)(nil)
	_ builder.CommandWithArgs        = (*Describe)(nil)
	_ builder.CommandWithFlags       = (*Describe)(nil)
	_ builder.CommandWithBasicClient = (*Describe)(nil)
	_ builder.CommandWithLogger      = (*Describe)(nil)
	_ builder.CommandWithExecute     = (*Describe)(nil)
)

type Describe struct {
	client global.BasicClient
	logger log.Logger
	args   struct {
		nameOrUUID string
	}
	flags struct {
		Path string `long:"path" usage:"Path to the app directory (default is local directory)"`
	}
}

func (d *Describe) Usage() string {
	return "describe [nameOrUUID] [--path pwd]"
}

func (d *Describe) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe a Conduit Data Application",
		Long: `This command will fetch details about the Application specified in '--path'
(or current working directory if not specified) on our Meroxa Platform,
or the Application specified by the given ID or Application Name.`,
		Example: `meroxa apps describe # assumes that the Application is in the current directory
meroxa apps describe --path /my/app
meroxa apps describe ID
meroxa apps describe NAME `,
	}
}

func (d *Describe) Execute(ctx context.Context) error {
	apps := &Applications{}
	var err error

	apps, err = RetrieveApplicationByNameOrID(ctx, d.client, d.args.nameOrUUID, d.flags.Path)
	if err != nil {
		return err
	}

	for _, app := range apps.Items {
		d.logger.Info(ctx, display.PrintTable(app, displayDetails))
		d.logger.JSON(ctx, app)
		dashboardURL := fmt.Sprintf("%s/conduitapps/%s/detail", global.GetMeroxaAPIURL(), app.ID)
		d.logger.Info(ctx, fmt.Sprintf("\n ✨ To view your application, visit %s", dashboardURL))
	}

	return nil
}

func (d *Describe) BasicClient(client global.BasicClient) {
	d.client = client
}

func (d *Describe) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Describe) ParseArgs(args []string) error {
	if len(args) > 0 {
		d.args.nameOrUUID = args[0]
	}

	return nil
}
