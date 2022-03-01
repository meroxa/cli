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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

type updateConnectorClient interface {
	UpdateConnectorStatus(ctx context.Context, nameOrID string, state meroxa.Action) (*meroxa.Connector, error)
	UpdateConnector(ctx context.Context, nameOrID string, input *meroxa.UpdateConnectorInput) (*meroxa.Connector, error)
}

type Update struct {
	client updateConnectorClient
	logger log.Logger

	args struct {
		NameOrID string
	}

	flags struct {
		Config string `long:"config" short:"c" usage:"new connector configuration"`
		Name   string `long:"name" usage:"new connector name"`
		State  string `long:"state" usage:"new connector state (pause | resume | restart)"`
	}
}

func (u *Update) Usage() string {
	return "update NAME"
}

func (u *Update) Docs() builder.Docs {
	return builder.Docs{
		Short: "Update connector name, configuration or state",
		Example: "\n" +
			"meroxa connector update old-name --name new-name' \n" +
			"meroxa connector update connector-name --state pause' \n" +
			"meroxa connector update connector-name --config '{\"table.name.format\":\"public.copy\"}' \n" +
			"meroxa connector update connector-name --state restart' \n",
	}
}

func (u *Update) Execute(ctx context.Context) error {
	// TODO: Implement something like dependent flags in Builder
	if u.flags.Config == "" && u.flags.Name == "" && u.flags.State == "" {
		return errors.New("requires either --config, --name or --state")
	}

	u.logger.Infof(ctx, "Updating connector %q...", u.args.NameOrID)
	var con *meroxa.Connector
	var err error

	if u.flags.State != "" {
		con, err = u.client.UpdateConnectorStatus(ctx, u.args.NameOrID, meroxa.Action(u.flags.State))
		if err != nil {
			return err
		}
	}

	if u.flags.Name != "" || u.flags.Config != "" {
		cu := &meroxa.UpdateConnectorInput{}

		// wants to update name
		if u.flags.Name != "" {
			cu.Name = u.flags.Name
		}

		// wants to update configuration
		if u.flags.Config != "" {
			config := map[string]interface{}{}

			err = json.Unmarshal([]byte(u.flags.Config), &config)
			if err != nil {
				return fmt.Errorf("can't parse config, make sure it is a valid JSON map: %w", err)
			}

			cu.Configuration = config
		}

		con, err = u.client.UpdateConnector(ctx, u.args.NameOrID, cu)
		if err != nil {
			return err
		}
	}

	u.logger.Infof(ctx, "Connector %q successfully updated!", u.args.NameOrID)
	u.logger.JSON(ctx, con)
	return nil
}

func (u *Update) Flags() []builder.Flag {
	return builder.BuildFlags(&u.flags)
}

func (u *Update) Logger(logger log.Logger) {
	u.logger = logger
}

func (u *Update) Client(client meroxa.Client) {
	u.client = client
}

func (u *Update) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires connector name")
	}

	u.args.NameOrID = args[0]
	return nil
}

var (
	_ builder.CommandWithDocs    = (*Update)(nil)
	_ builder.CommandWithArgs    = (*Update)(nil)
	_ builder.CommandWithFlags   = (*Update)(nil)
	_ builder.CommandWithClient  = (*Update)(nil)
	_ builder.CommandWithLogger  = (*Update)(nil)
	_ builder.CommandWithExecute = (*Update)(nil)
)
