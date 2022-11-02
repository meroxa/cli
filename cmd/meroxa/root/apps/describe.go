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
	"errors"

	"github.com/meroxa/cli/utils/display"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Describe)(nil)
	_ builder.CommandWithArgs    = (*Describe)(nil)
	_ builder.CommandWithClient  = (*Describe)(nil)
	_ builder.CommandWithLogger  = (*Describe)(nil)
	_ builder.CommandWithExecute = (*Describe)(nil)
)

type describeApplicationClient interface {
	GetApplication(ctx context.Context, nameOrUUID string) (*meroxa.Application, error)
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
	GetConnectorByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Connector, error)
	GetFunction(ctx context.Context, nameOrUUID string) (*meroxa.Function, error)
}

type Describe struct {
	client describeApplicationClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}
}

func (d *Describe) Usage() string {
	return "describe [NAMEorUUID]"
}

func (d *Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe a Turbine Data Application",
		Beta:  true,
	}
}

func (d *Describe) Execute(ctx context.Context) error {
	var output string

	app, err := d.client.GetApplication(ctx, d.args.NameOrUUID)
	if err != nil {
		return err
	}

	output = display.AppTable(app)

	d.logger.Info(ctx, output)
	d.logger.JSON(ctx, app)

	return nil
}

func (d *Describe) Client(client meroxa.Client) {
	d.client = client
}

func (d *Describe) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Describe) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires app name or UUID")
	}

	d.args.NameOrUUID = args[0]
	return nil
}
