/*
Copyright © 2023 Meroxa Inc

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

package introspect

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

type refreshResourceIntrospectionClient interface {
	// RefreshIntrospection
}

type Refresh struct {
	client refreshResourceIntrospectionClient
	logger log.Logger

	args struct {
		ResourceNameOrUUID string
	}
}

var (
	_ builder.CommandWithDocs    = (*Refresh)(nil)
	_ builder.CommandWithArgs    = (*Refresh)(nil)
	_ builder.CommandWithClient  = (*Refresh)(nil)
	_ builder.CommandWithLogger  = (*Refresh)(nil)
	_ builder.CommandWithExecute = (*Refresh)(nil)
)

func (c *Refresh) Usage() string {
	return "create [ResourceNameOrUUID]"
}

func (c *Refresh) Docs() builder.Docs {
	return builder.Docs{
		Short: "Add a resource to your Meroxa resource catalog",
		Long:  `Use the create command to add resources to your Meroxa resource catalog.`,

		// TODO: Provide example with `--env` once it's not behind a feature flag
		Example: `
$ meroxa resource create mybigquery \
    --type bigquery \
    -u "bigquery://$GCP_PROJECT_ID/$GCP_DATASET_NAME" \
    --client-key "$(cat $GCP_SERVICE_ACCOUNT_JSON_FILE)"
`,
	}
}

func (c *Refresh) Client(client meroxa.Client) {
	c.client = client
}

func (c *Refresh) Logger(logger log.Logger) {
	c.logger = logger
}

func (c *Refresh) ParseArgs(args []string) error {
	if len(args) > 0 {
		c.args.ResourceNameOrUUID = args[0]
	}
	return nil
}

func (c *Refresh) Execute(ctx context.Context) error {
	return nil
}
