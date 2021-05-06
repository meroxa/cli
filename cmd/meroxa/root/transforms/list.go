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

package transforms

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

var (
	_ builder.CommandWithDocs    = (*List)(nil)
	_ builder.CommandWithClient  = (*List)(nil)
	_ builder.CommandWithLogger  = (*List)(nil)
	_ builder.CommandWithExecute = (*List)(nil)
	_ builder.CommandWithAliases = (*List)(nil)
)

type listTransformsClient interface {
	ListTransforms(ctx context.Context) ([]*meroxa.Transform, error)
}

type List struct {
	client listTransformsClient
	logger log.Logger
}

func (l *List) Usage() string {
	return "list"
}

func (l *List) Docs() builder.Docs {
	return builder.Docs{
		Short: "List transforms",
	}
}

func (l *List) Aliases() []string {
	return []string{"ls"}
}

func (l *List) Execute(ctx context.Context) error {
	var err error
	transforms, err := l.client.ListTransforms(ctx)
	if err != nil {
		return err
	}

	l.logger.JSON(ctx, transforms)
	l.logger.Info(ctx, utils.TransformsTable(transforms))

	return nil
}

func (l *List) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *List) Client(client *meroxa.Client) {
	l.client = client
}