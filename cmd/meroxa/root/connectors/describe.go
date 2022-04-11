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
	"context"
	"errors"

	"github.com/meroxa/cli/utils"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs       = (*Describe)(nil)
	_ builder.CommandWithArgs       = (*Describe)(nil)
	_ builder.CommandWithClient     = (*Describe)(nil)
	_ builder.CommandWithLogger     = (*Describe)(nil)
	_ builder.CommandWithExecute    = (*Describe)(nil)
	_ builder.CommandWithDeprecated = (*Describe)(nil)
)

type describeConnectorClient interface {
	GetConnectorByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Connector, error)
}

type Describe struct {
	client describeConnectorClient
	logger log.Logger

	args struct {
		NameOrID string
	}
}

func (d *Describe) Usage() string {
	return "describe [NAME]"
}

func (d *Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe connector",
	}
}

func (d *Describe) Execute(ctx context.Context) error {
	connector, err := d.client.GetConnectorByNameOrID(ctx, d.args.NameOrID)
	if err != nil {
		return err
	}

	d.logger.Info(ctx, utils.ConnectorTable(connector))
	d.logger.JSON(ctx, connector)

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
		return errors.New("requires connector name")
	}

	d.args.NameOrID = args[0]
	return nil
}

func (*Describe) Deprecated() string {
	return "we encourage you to describe your application via `meroxa apps describe` instead."
}
