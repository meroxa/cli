/*
Copyright © 2020 Meroxa Inc

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

package add

import (
	"context"
	"encoding/json"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
)

type addResourceClient interface {
	CreateResource(ctx context.Context, resource *meroxa.CreateResourceInput) (*meroxa.Resource, error)
}

type AddResource struct {
	client addResourceClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Type        string `long:"type"        short:""  usage:"resource type"        required:"true"`
		Url         string `long:"url"         short:"u" usage:"resource url"         required:"true"`
		Metadata    string `long:"credentials" short:""  usage:"resource credentials" required:"false"`
		Credentials string `long:"metadata"    short:"m" usage:"resource metadata"    required:"false"`
	}
}

var (
	_ builder.CommandWithDocs    = (*AddResource)(nil)
	_ builder.CommandWithArgs    = (*AddResource)(nil)
	_ builder.CommandWithFlags   = (*AddResource)(nil)
	_ builder.CommandWithClient  = (*AddResource)(nil)
	_ builder.CommandWithLogger  = (*AddResource)(nil)
	_ builder.CommandWithExecute = (*AddResource)(nil)
)

func (ar *AddResource) Usage() string {
	return "resource [NAME] --type TYPE --url URL"
}

func (ar *AddResource) Docs() builder.Docs {
	return builder.Docs{
		Short: "Add a resource to your Meroxa resource catalog",
		Long:  `Use the add command to add resources to your Meroxa resource catalog.`,
		Example: "\n" +
			"meroxa add resource store --type postgres -u $DATABASE_URL --metadata '{\"logical_replication\":true}'\n" +
			"meroxa add resource datalake --type s3 -u \"s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos\"\n" +
			"meroxa add resource warehouse --type redshift -u $REDSHIFT_URL\n" +
			"meroxa add resource slack --type url -u $WEBHOOK_URL\n",
	}
}

func (ar *AddResource) Client(client *meroxa.Client) {
	ar.client = client
}

func (ar *AddResource) Logger(logger log.Logger) {
	ar.logger = logger
}

func (ar *AddResource) Flags() []builder.Flag {
	return builder.BuildFlags(&ar.flags)
}

func (ar *AddResource) ParseArgs(args []string) error {
	if len(args) > 0 {
		ar.args.Name = args[0]
	}
	return nil
}

func (ar *AddResource) Execute(ctx context.Context) error {
	input := meroxa.CreateResourceInput{
		Type:     ar.flags.Type,
		Name:     ar.args.Name,
		URL:      ar.flags.Url,
		Metadata: nil,
	}

	ar.logger.Infof(ctx, "Adding %s resource...", input.Type)

	// TODO: Figure out best way to handle creds and metadata
	// Get credentials (expect a JSON string)
	if ar.flags.Credentials != "" {
		err := json.Unmarshal([]byte(ar.flags.Credentials), &input.Credentials)
		if err != nil {
			return err
		}
	}

	if ar.flags.Metadata != "" {
		err := json.Unmarshal([]byte(ar.flags.Metadata), &input.Metadata)
		if err != nil {
			return err
		}
	}

	res, err := ar.client.CreateResource(ctx, &input)
	if err != nil {
		return err
	}

	ar.logger.Infof(ctx, "%s resource with name %s successfully added!", res.Type, res.Name)
	ar.logger.JSON(ctx, res)

	return nil
}
