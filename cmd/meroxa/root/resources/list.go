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

package resources

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
	_ builder.CommandWithAliases   = (*List)(nil)
	_ builder.CommandWithFlags     = (*List)(nil)
	_ builder.CommandWithNoHeaders = (*List)(nil)
)

type listResourcesClient interface {
	ListResources(ctx context.Context) ([]*meroxa.Resource, error)
	ListResourceTypes(ctx context.Context) ([]string, error)
}

type List struct {
	client      listResourcesClient
	logger      log.Logger
	hideHeaders bool

	flags struct {
		Types bool `long:"types" short:"" usage:"list resource types"`
		Type  bool `long:"type" short:"" usage:"alias to --types" hidden:"true"`
	}

	// ListTypes is used by the alias `meroxa list resource-types`.
	// Once we stop giving support to v1 commands, this could be removed
	ListTypes bool
}

func (l *List) Usage() string {
	return "list"
}

func (l *List) Docs() builder.Docs {
	return builder.Docs{
		Short: "List resources and resource types",
	}
}

func (l *List) Flags() []builder.Flag {
	return builder.BuildFlags(&l.flags)
}

func (l *List) Aliases() []string {
	return []string{"ls"}
}

func (l *List) Execute(ctx context.Context) error {
	var err error

	// What used to be `meroxa list resource-types`
	if l.flags.Types || l.flags.Type || l.ListTypes {
		var rTypes []string

		rTypes, err = l.client.ListResourceTypes(ctx)

		if err != nil {
			return err
		}

		l.logger.JSON(ctx, rTypes)
		l.logger.Info(ctx, utils.ResourceTypesTable(rTypes, l.hideHeaders))

		return nil
	}

	resources, err := l.client.ListResources(ctx)
	if err != nil {
		return err
	}

	l.logger.JSON(ctx, resources)
	l.logger.Info(ctx, utils.ResourcesTable(resources, l.hideHeaders))

	return nil
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
