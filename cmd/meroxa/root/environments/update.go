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

package environments

import (
	"context"
	"fmt"

	"github.com/meroxa/cli/utils"
	"github.com/meroxa/cli/utils/display"

	"github.com/manifoldco/promptui"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Update)(nil)
	_ builder.CommandWithArgs    = (*Update)(nil)
	_ builder.CommandWithFlags   = (*Update)(nil)
	_ builder.CommandWithClient  = (*Update)(nil)
	_ builder.CommandWithLogger  = (*Update)(nil)
	_ builder.CommandWithExecute = (*Update)(nil)
	_ builder.CommandWithPrompt  = (*Update)(nil)
)

type updateEnvironmentClient interface {
	UpdateEnvironment(ctx context.Context, nameOrUUID string, body *meroxa.UpdateEnvironmentInput) (*meroxa.Environment, error)
}

type Update struct {
	client updateEnvironmentClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}

	flags struct {
		Name   string   `long:"name" usage:"updated environment name, when specified"`
		Config []string `short:"c" long:"config" usage:"updated environment configuration based on type and provider (e.g.: --config aws_access_key_id=my_access_key --config aws_secret_access_key=my_access_secret)"` //nolint:lll
	}

	envCfg map[string]interface{}
}

func (c *Update) Logger(logger log.Logger) {
	c.logger = logger
}

func (c *Update) Client(client meroxa.Client) {
	c.client = client
}

func (c *Update) Flags() []builder.Flag {
	return builder.BuildFlags(&c.flags)
}

func (c *Update) ParseArgs(args []string) error {
	if len(args) > 0 {
		c.args.NameOrUUID = args[0]
	}
	return nil
}

func (c *Update) SkipPrompt() bool {
	return c.args.NameOrUUID != "" && c.flags.Name != "" && len(c.flags.Config) != 0
}

func (c *Update) setUserValues(e *meroxa.UpdateEnvironmentInput) {
	if c.flags.Name != "" {
		e.Name = c.flags.Name
	}

	if c.flags.Config != nil {
		e.Configuration = utils.StringSliceToInterfaceMap(c.flags.Config)
	}
}

func (c *Update) Execute(ctx context.Context) error {
	e := &meroxa.UpdateEnvironmentInput{}
	c.setUserValues(e)

	// In case user skipped prompt and configuration was specified via flags
	if len(c.flags.Config) != 0 {
		e.Configuration = utils.StringSliceToInterfaceMap(c.flags.Config)
	}
	if c.flags.Name != "" {
		e.Name = c.flags.Name
	}

	c.logger.Infof(ctx, "Updating environment...")

	environment, err := c.client.UpdateEnvironment(ctx, c.args.NameOrUUID, e)
	if err != nil {
		return err
	}

	if environment.Status.State != meroxa.EnvironmentStatePreflightSuccess {
		details := display.EnvironmentPreflightTable(environment)
		c.logger.Errorf(ctx,
			"Environment %q could not be updated because it failed the preflight checks\n%s\n",
			environment.Name,
			details)
	} else {
		c.logger.Infof(ctx,
			"Preflight checks have passed. Environment %q is being updated. Run `meroxa env describe %s` for status",
			environment.Name,
			environment.Name)
	}

	c.logger.JSON(ctx, environment)
	return nil
}

func (c *Update) NotConfirmed() string {
	return "\nTo view all options for updating a Meroxa Environment,\n " +
		"please run \"meroxa help env update\". \n"
}

func (c *Update) showEventConfirmation() {
	var eventToConfirm string

	eventToConfirm = "Environment details:\n"
	eventToConfirm += fmt.Sprintf("\tCurrent Name or UUID: %s\n", c.args.NameOrUUID)
	if c.flags.Name != "" {
		eventToConfirm += fmt.Sprintf("\tNew Name: %s\n", c.flags.Name)
	}

	if len(c.envCfg) > 0 {
		eventToConfirm += "\tNew Config:"

		for k, v := range c.envCfg {
			eventToConfirm += fmt.Sprintf("\n\t  %s: %s", k, v)
		}
	}

	fmt.Println(eventToConfirm)
}

func (c *Update) Prompt() error {
	if c.args.NameOrUUID == "" {
		p := promptui.Prompt{
			Label:   "Current Environment name or UUID",
			Default: "",
		}

		c.args.NameOrUUID, _ = p.Run()
	}

	if c.flags.Name == "" {
		p := promptui.Prompt{
			Label:   "New Environment name (optional)",
			Default: "",
		}

		c.flags.Name, _ = p.Run()
	}

	c.envCfg = utils.StringSliceToInterfaceMap(c.flags.Config)
	configPrompt := "a"
	if len(c.flags.Config) != 0 {
		configPrompt = "additional"
	}

	p := promptui.Prompt{
		Label:     fmt.Sprintf("Does your environment require %s new configuration", configPrompt),
		IsConfirm: true,
	}

	_, err := p.Run()

	// user responded "yes" to confirmation prompt
	if err == nil {
		cfgIsNeeded := true

		for cfgIsNeeded {
			p = promptui.Prompt{
				Label: "Configuration key",
			}

			k, _ := p.Run()

			p = promptui.Prompt{
				Label: k,
			}

			v, _ := p.Run()
			c.envCfg[k] = v

			p := promptui.Prompt{
				Label:     "Add another configuration",
				IsConfirm: true,
			}

			_, err = p.Run()
			if err != nil {
				cfgIsNeeded = false
			}
		}
	}

	c.showEventConfirmation()

	prompt := promptui.Prompt{
		Label:     "Update this environment",
		IsConfirm: true,
	}

	_, err = prompt.Run()
	return err
}

func (c *Update) Usage() string {
	return "update NAMEorUUID"
}

func (c *Update) Docs() builder.Docs {
	return builder.Docs{
		Short: "Update an environment",
		Example: `
meroxa env update my-env --name new-name --config aws_access_key_id=my_access_key --config aws_access_secret=my_access_secret"
`,
	}
}
