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

package apps

import (
	"context"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Create)(nil)
	_ builder.CommandWithArgs    = (*Create)(nil)
	_ builder.CommandWithFlags   = (*Create)(nil)
	_ builder.CommandWithClient  = (*Create)(nil)
	_ builder.CommandWithLogger  = (*Create)(nil)
	_ builder.CommandWithExecute = (*Create)(nil)
)

type createApplicationClient interface {
	CreateApplication(ctx context.Context, input *meroxa.CreateApplicationInput) (*meroxa.Application, error)
}

type Create struct {
	client createApplicationClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Path string `long:"path" usage:"path where application will be initialized (current directory as default)"`
	}
}

func (c *Create) Logger(logger log.Logger) {
	c.logger = logger
}

func (c *Create) Client(client meroxa.Client) {
	c.client = client
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

func (c *Create) Execute(ctx context.Context) error {
	path := turbineCLI.GetPath(c.flags.Path)
	lang, err := turbineCLI.GetLangFromAppJSON(path)
	if err != nil {
		return err
	}

	input := meroxa.CreateApplicationInput{
		Name:     c.args.Name,
		Language: lang,
	}

	c.logger.Infof(ctx, "Creating application %q...", input.Name)

	res, err := c.client.CreateApplication(ctx, &input)
	if err != nil {
		return err
	}

	c.logger.Infof(ctx, "Application %q successfully created!", res.Name)
	c.logger.JSON(ctx, res)

	return nil
}

func (c *Create) Usage() string {
	return "create NAME"
}

func (c *Create) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Create an application",
		Example: "meroxa apps create my-app --language golang",
	}
}
