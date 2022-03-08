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
	"errors"
	"fmt"

	"github.com/volatiletech/null/v8"

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
		Path     string `long:"path" usage:"path where app was initialized to read its language from app.json (required unless specifying --lang)"` //nolint:lll
		Lang     string `long:"language" short:"l" usage:"language to use (required unless specifying --path)"`
		Pipeline string `long:"pipeline" short:"p" usage:"UUID of pipeline associated with this application"`
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
	if len(args) < 1 {
		return errors.New("requires an application name")
	}
	c.args.Name = args[0]
	return nil
}

func (c *Create) getLang(ctx context.Context) (string, error) {
	if path := c.flags.Path; path != "" {
		lang, err := turbineCLI.GetLangFromAppJSON(path)
		if err != nil {
			return lang, err
		}
		if c.flags.Lang != "" && c.flags.Lang != lang {
			c.logger.Info(ctx, "Ignoring language flag.")
		}
		return lang, nil
	}

	return c.flags.Lang, nil
}

func (c *Create) Execute(ctx context.Context) error {
	if c.flags.Lang == "" && c.flags.Path == "" {
		return fmt.Errorf("language is required either using --path ~/turbine/my-app or --lang. Type `meroxa help apps create` for more information") //nolint:lll
	}

	lang, err := c.getLang(ctx)
	if err != nil {
		return err
	}

	input := meroxa.CreateApplicationInput{
		Name:     c.args.Name,
		Language: lang,
		Pipeline: meroxa.EntityIdentifier{UUID: null.StringFrom(c.flags.Pipeline)},
	}

	c.logger.Infof(ctx, "Creating application %q with language %q...", input.Name, lang)

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

//nolint:lll
func (c *Create) Docs() builder.Docs {
	return builder.Docs{
		Short: "Create a Meroxa Data Application",
		Long: "You'll be able to use this application for consequent build via `meroxa apps deploy`. You'll need to specify " +
			"language used either via `--language` or specifying with `--path` the location of your app.json which should contain the desired language.",
		Example: "meroxa apps create my-app --language golang\n" +
			"meroxa apps create my-app --path ~/turbine/my-app",
	}
}
