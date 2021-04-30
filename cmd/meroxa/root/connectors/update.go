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

package connectors

import (
	"context"
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/meroxa-go"
)

type updateConnectorClient interface {
	UpdateConnectorStatus(ctx context.Context, connectorKey, state string) (*meroxa.Connector, error)
}

type UpdateConnector struct {
	client updateConnectorClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		State string `long:"state" usage:"connector state" required:"true"`
	}
}

func (u *UpdateConnector) Usage() string {
	return "update NAME --state pause | resume | restart"
}

func (u *UpdateConnector) Docs() builder.Docs {
	return builder.Docs{
		Short: "Update connector state",
	}
}

func (u *UpdateConnector) Execute(ctx context.Context) error {
	u.logger.Infof(ctx, "Updating connector %q...", u.args.Name)

	con, err := u.client.UpdateConnectorStatus(ctx, u.args.Name, u.flags.State)
	if err != nil {
		return err
	}

	u.logger.Infof(ctx, "Connector %q successfully updated!", u.args.Name)
	u.logger.JSON(ctx, con)
	return nil
}

func (u *UpdateConnector) Flags() []builder.Flag {
	return builder.BuildFlags(&u.flags)
}

func (u *UpdateConnector) Logger(logger log.Logger) {
	u.logger = logger
}

func (u *UpdateConnector) Client(client *meroxa.Client) {
	u.client = client
}

func (u *UpdateConnector) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires connector name")
	}

	u.args.Name = args[0]
	return nil
}

var (
	_ builder.CommandWithDocs    = (*UpdateConnector)(nil)
	_ builder.CommandWithArgs    = (*UpdateConnector)(nil)
	_ builder.CommandWithFlags   = (*UpdateConnector)(nil)
	_ builder.CommandWithClient  = (*UpdateConnector)(nil)
	_ builder.CommandWithLogger  = (*UpdateConnector)(nil)
	_ builder.CommandWithExecute = (*UpdateConnector)(nil)
)
