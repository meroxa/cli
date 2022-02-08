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

package environments

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Create)(nil)
	_ builder.CommandWithArgs    = (*Create)(nil)
	_ builder.CommandWithFlags   = (*Create)(nil)
	_ builder.CommandWithClient  = (*Create)(nil)
	_ builder.CommandWithLogger  = (*Create)(nil)
	_ builder.CommandWithExecute = (*Create)(nil)
	_ builder.CommandWithPrompt  = (*Create)(nil)
)

type createEnvironmentClient interface {
	CreateEnvironment(ctx context.Context, body *meroxa.CreateEnvironmentInput) (*meroxa.Environment, error)
}

type Create struct {
	client createEnvironmentClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Type     string   `long:"type" usage:"environment type, when not specified"`
		Provider string   `long:"provider" usage:"environment cloud provider to use"`
		Region   string   `long:"region" usage:"environment region"`
		Config   []string `short:"c" long:"config" usage:"environment configuration based on type and provider (e.g.: --config aws_access_key_id=my_access_key --config aws_secret_access_key=my_access_secret)"` // nolint:lll
	}

	envCfg map[string]interface{}
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

func (c *Create) SkipPrompt() bool {
	return c.args.Name != "" && c.flags.Type != "" && c.flags.Provider != "" && c.flags.Region != "" && len(c.flags.Config) != 0
}

func (c *Create) setUserValues(e *meroxa.CreateEnvironmentInput) {
	if c.args.Name != "" {
		e.Name = c.args.Name
	}

	if c.flags.Type != "" {
		e.Type = meroxa.EnvironmentType(c.flags.Type)
	}

	if c.flags.Provider != "" {
		e.Provider = meroxa.EnvironmentProvider(c.flags.Provider)
	}

	if c.flags.Region != "" {
		e.Region = meroxa.EnvironmentRegion(c.flags.Region)
	}

	if c.envCfg != nil {
		e.Configuration = c.envCfg
	}
}

func stringSliceToMap(input []string) map[string]interface{} {
	const pair = 2
	m := make(map[string]interface{})
	for _, config := range input {
		parts := strings.Split(config, "=")
		if len(parts) >= pair {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func (c *Create) Execute(ctx context.Context) error {
	e := &meroxa.CreateEnvironmentInput{}
	c.setUserValues(e)

	// In case user skipped prompt and configuration was specified via flags
	if len(c.flags.Config) != 0 {
		e.Configuration = stringSliceToMap(c.flags.Config)
	}

	c.logger.Infof(ctx, "Provisioning environment...")

	environment, err := c.client.CreateEnvironment(ctx, e)

	if err != nil {
		return err
	}

	if environment.Status.State == meroxa.EnvironmentStatePreflightError {
		details := utils.PrettyString(environment.Status.PreflightDetails)
		c.logger.Errorf(ctx, "Environment %q could not be provisioned because it failed the preflight checks\n%s\n", environment.Name, details)
	} else {
		environment.Status.PreflightDetails = nil
		c.logger.Infof(ctx, "Preflight checks have passed. Environment %q is being provisioned. Run `meroxa env describe %s` for status", environment.Name, environment.Name)
	}

	c.logger.JSON(ctx, environment)
	return nil
}

func (c *Create) NotConfirmed() string {
	return "\nTo view all options for creating a Meroxa Environment,\n " +
		"please run \"meroxa help env create\". \n"
}

func (c *Create) showEventConfirmation() {
	var eventToConfirm string

	eventToConfirm = "Environment details:\n"

	if c.args.Name != "" {
		eventToConfirm += fmt.Sprintf("\tName: %s\n", c.args.Name)
	}

	eventToConfirm += fmt.Sprintf("\tType: %s\n\tProvider: %s\n\tRegion: %s\n", c.flags.Type, c.flags.Provider, c.flags.Region)

	if len(c.envCfg) > 0 {
		eventToConfirm += "\tConfig:"

		for k, v := range c.envCfg {
			eventToConfirm += fmt.Sprintf("\n\t  %s: %s", k, v)
		}
	}

	fmt.Println(eventToConfirm)
}

func (c *Create) Prompt() error {
	if c.args.Name == "" {
		p := promptui.Prompt{
			Label:   "Environment name (optional)",
			Default: "",
		}

		c.args.Name, _ = p.Run()
	}

	if c.flags.Type == "" {
		vType := func(input string) error {
			switch input {
			case "self_hosted", "private":
				return nil
			default:
				return errors.New("unsupported environment type")
			}
		}

		p := promptui.Prompt{
			Label:    "Type (self_hosted or private)",
			Default:  "self_hosted",
			Validate: vType,
		}

		c.flags.Type, _ = p.Run()
	}

	if c.flags.Provider == "" {
		p := promptui.Prompt{
			Label:   "Cloud provider",
			Default: "aws",
		}

		c.flags.Provider, _ = p.Run()
	}

	if c.flags.Region == "" {
		p := promptui.Prompt{
			Label:   "Region",
			Default: "us-east-1",
		}

		c.flags.Region, _ = p.Run()
	}

	if len(c.flags.Config) != 0 {
		c.envCfg = make(map[string]interface{})

		p := promptui.Prompt{
			Label:     "Does your environment require configuration",
			IsConfirm: true,
		}

		_, err := p.Run()

		// user responded "yes" to confirmation prompt
		if err == nil {
			cfgIsNeeded := true

			for cfgIsNeeded {
				p = promptui.Prompt{
					Label: "\tConfiguration key",
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

				_, err := p.Run()
				if err != nil {
					cfgIsNeeded = false
				}
			}
		}
	}

	c.showEventConfirmation()

	prompt := promptui.Prompt{
		Label:     "Create this environment",
		IsConfirm: true,
	}

	_, err := prompt.Run()
	return err
}

func (c *Create) Usage() string {
	return "create NAME"
}

//nolint:lll
func (c *Create) Docs() builder.Docs {
	return builder.Docs{
		Short: "Create an environment",
		Example: `
meroxa env create my-env --type self_hosted --provider aws --region us-east-1 --config aws_access_key_id=my_access_key --config aws_secret_access_key=my_access_secret
`,
	}
}
