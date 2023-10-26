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
	"encoding/json"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
)

var (
	_ builder.CommandWithDocs        = (*List)(nil)
	_ builder.CommandWithBasicClient = (*List)(nil)
	_ builder.CommandWithLogger      = (*List)(nil)
	_ builder.CommandWithExecute     = (*List)(nil)
	_ builder.CommandWithAliases     = (*List)(nil)
)

type List struct {
	client global.BasicClient
	logger log.Logger
}

func (l *List) Usage() string {
	return "list"
}

func (l *List) Docs() builder.Docs {
	return builder.Docs{
		Short: "List Turbine Data Applications",
	}
}

func (l *List) Aliases() []string {
	return []string{"ls"}
}

func (l *List) Execute(ctx context.Context) error {
	var err error
	apps := &Applications{}

	response, err := l.client.CollectionRequest(ctx, "GET", collectionName, "", nil, nil, apps)
	if err != nil {
		return err
	}
	err = json.NewDecoder(response.Body).Decode(&apps)
	if err != nil {
		return err
	}

	l.logger.Info(ctx, display.PrintList(apps.Items, displayDetails))
	l.logger.JSON(ctx, apps)

	output := " ✨ To view your applications, visit https://dashboard.meroxa.io/apps"
	l.logger.Info(ctx, output)
	return nil
}

func (l *List) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *List) BasicClient(client global.BasicClient) {
	l.client = client
}
