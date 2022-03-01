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

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs      = (*List)(nil)
	_ builder.CommandWithClient    = (*List)(nil)
	_ builder.CommandWithLogger    = (*List)(nil)
	_ builder.CommandWithExecute   = (*List)(nil)
	_ builder.CommandWithFlags     = (*List)(nil)
	_ builder.CommandWithAliases   = (*List)(nil)
	_ builder.CommandWithNoHeaders = (*List)(nil)
)

type listConnectorsClient interface {
	ListConnectors(ctx context.Context) ([]*meroxa.Connector, error)
	ListPipelineConnectors(ctx context.Context, pipelineID int) ([]*meroxa.Connector, error)
	GetPipelineByName(ctx context.Context, name string) (*meroxa.Pipeline, error)
}

type List struct {
	client      listConnectorsClient
	logger      log.Logger
	hideHeaders bool

	flags struct {
		Pipeline string `long:"pipeline" short:""  usage:"filter connectors by pipeline name"`
	}
}

func (l *List) Usage() string {
	return "list"
}

func (l *List) Docs() builder.Docs {
	return builder.Docs{
		Short: "List connectors",
	}
}

func (l *List) Aliases() []string {
	return []string{"ls"}
}

func (l *List) Execute(ctx context.Context) error {
	var err error
	var connectors []*meroxa.Connector

	// Filtering by pipeline name
	if l.flags.Pipeline != "" {
		var p *meroxa.Pipeline

		p, err = l.client.GetPipelineByName(ctx, l.flags.Pipeline)

		if err != nil {
			return err
		}

		connectors, err = l.client.ListPipelineConnectors(ctx, p.ID)

		if err != nil {
			return err
		}
	} else {
		connectors, err = l.client.ListConnectors(ctx)

		if err != nil {
			return err
		}
	}

	l.logger.JSON(ctx, connectors)
	l.logger.Info(ctx, utils.ConnectorsTable(connectors, l.hideHeaders))

	return nil
}

func (l *List) Flags() []builder.Flag {
	return builder.BuildFlags(&l.flags)
}

func (l *List) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *List) Client(client meroxa.Client) {
	l.client = client
}

func (l *List) HideHeaders(hide bool) {
	l.hideHeaders = hide
}
