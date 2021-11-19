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

package endpoints

import (
	"context"

	"github.com/meroxa/meroxa-go/pkg/meroxa"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type createEndpointClient interface {
	CreateEndpoint(ctx context.Context, input *meroxa.CreateEndpointInput) error
}

type Create struct {
	client createEndpointClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Protocol string `long:"protocol" short:"p" usage:"protocol, value can be http or grpc" required:"true"`
		Stream   string `long:"stream" short:"s" usage:"stream name" required:"true"`
	}
}

func (c *Create) Usage() string {
	return "create [NAME] [flags]"
}

func (c *Create) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Create an endpoint",
		Long:    "Use create endpoint to expose an endpoint to a connector stream",
		Example: "meroxa endpoints create my-endpoint --protocol http --stream my-stream",
	}
}

func (c *Create) Execute(ctx context.Context) error {
	c.logger.Info(ctx, "Creating endpoint...")
	input := &meroxa.CreateEndpointInput{
		Name:     c.args.Name,
		Protocol: meroxa.EndpointProtocol(c.flags.Protocol),
		Stream:   c.flags.Stream,
	}

	err := c.client.CreateEndpoint(ctx, input)

	if err != nil {
		return err
	}

	c.logger.Info(ctx, "Endpoint successfully created!")

	return nil
}

func (c *Create) Client(client meroxa.Client) {
	c.client = client
}

func (c *Create) Logger(logger log.Logger) {
	c.logger = logger
}

func (c *Create) Flags() []builder.Flag {
	return builder.BuildFlags(&c.flags)
}

func (c *Create) ParseArgs(args []string) error {
	if len(args) > 0 {
		c.args.Name = args[0]
	}
	return nil
}

var (
	_ builder.CommandWithDocs    = (*Create)(nil)
	_ builder.CommandWithArgs    = (*Create)(nil)
	_ builder.CommandWithFlags   = (*Create)(nil)
	_ builder.CommandWithClient  = (*Create)(nil)
	_ builder.CommandWithLogger  = (*Create)(nil)
	_ builder.CommandWithExecute = (*Create)(nil)
)
