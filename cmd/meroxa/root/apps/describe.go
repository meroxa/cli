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

package apps

import (
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
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
	GetFunctionLogs(ctx context.Context, nameOrUUID string) (*http.Response, error)
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

	flags struct {
		Extended bool `long:"extended" usage:"whether to show additional details about the Turbine Data Application"`
	}
}

func (d *Describe) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Describe) Usage() string {
	return "describe [NAMEorUUID]"
}

func (d *Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe a Turbine Data Application",
	}
}

func (d *Describe) Execute(ctx context.Context) error {
	var output string

	app, err := d.client.GetApplication(ctx, d.args.NameOrUUID)
	if err != nil {
		return err
	}

	if d.flags.Extended {
		resources := make([]*meroxa.Resource, 0)
		connectors := make(map[string]*meroxa.Connector)
		functions := make([]*meroxa.Function, 0)

		for _, rr := range app.Resources {
			resource, err := d.client.GetResourceByNameOrID(ctx, rr.Name.String)
			if err != nil {
				return err
			}
			resources = append(resources, resource)
		}
		for _, cc := range app.Connectors {
			connector, err := d.client.GetConnectorByNameOrID(ctx, cc.Name.String)
			if err != nil {
				return err
			}
			connectors[connector.ResourceName] = connector
		}
		for _, ff := range app.Functions {
			function, err := d.client.GetFunction(ctx, ff.UUID.String)
			if err != nil {
				return err
			}

			// Include logs
			resp, err := d.client.GetFunctionLogs(ctx, ff.Name.String)
			if err != nil {
				return err
			}
			defer resp.Body.Close() //nolint:gocritic

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			if err != nil {
				return err
			}

			function.Logs = buf.String()

			functions = append(functions, function)
		}

		output = utils.AppExtendedTable(app, resources, connectors, functions)
	} else {
		output = utils.AppTable(app)
	}

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
		return errors.New("requires app name")
	}

	d.args.NameOrUUID = args[0]
	return nil
}
