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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/prompt"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
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

const (
	AwsProvider = "aws"
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
		Type     string `long:"type" usage:"environment type, when not specified"`
		Provider string `long:"provider" usage:"environment cloud provider to use"`
		Region   string `long:"region" usage:"environment region"`
		Config   string `short:"c" long:"config" usage:"environment configuration based on type and provider (e.g.: --config {\"aws_access_key_id\":\"my_access_key\", \"aws_access_secret\":\"my_access_secret\"})"` // nolint:lll
	}

	envCfg map[interface{}]interface{}
}

func (c *Create) Logger(logger log.Logger) {
	c.logger = logger
}

func (c *Create) Client(client *meroxa.Client) {
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

func (c *Create) setUserValues(e *meroxa.CreateEnvironmentInput) {
	if c.args.Name != "" {
		e.Name = c.args.Name
	}

	if c.flags.Type != "" {
		e.Type = c.flags.Type
	}

	if c.flags.Provider != "" {
		e.Provider = c.flags.Provider
	}

	if c.flags.Region != "" {
		e.Region = c.flags.Region
	}

	if c.envCfg != nil {
		envCfg := make(map[string]interface{})
		for k, v := range c.envCfg {
			envCfg[k.(string)] = v
		}
		e.Configuration = envCfg
	}
}

func (c *Create) Execute(ctx context.Context) error {
	e := &meroxa.CreateEnvironmentInput{}
	c.setUserValues(e)

	// In case user skipped prompt and configuration was specified via flags
	if c.flags.Config != "" {
		var config map[string]interface{}
		err := json.Unmarshal([]byte(c.flags.Config), &config)
		if err != nil {
			return fmt.Errorf("could not parse configuration: %w", err)
		}

		e.Configuration = config
	}

	c.logger.Infof(ctx, "Provisioning environment...")

	environment, err := c.client.CreateEnvironment(ctx, e)

	if err != nil {
		return err
	}

	c.logger.Infof(ctx, "Environment %q is being provisioned. Run `meroxa env describe %s` for status", environment.Name, environment.Name)
	c.logger.JSON(ctx, environment)

	return nil
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

func (c *Create) Prompt(ctx context.Context) error {
	c.envCfg = make(map[interface{}]interface{})

	err := prompt.Show(ctx, []prompt.Prompt{
		prompt.StringPrompt{
			Label: "Environment name (optional)",
			Value: &c.args.Name,
			Skip:  c.args.Name != "",
		},
		prompt.StringPrompt{
			Label:   "Type (hosted or dedicated)",
			Default: "hosted",
			Validate: func(input string) error {
				switch input {
				case "hosted", "dedicated":
					return nil
				default:
					return errors.New("unsupported environment type")
				}
			},
			Value: &c.flags.Type,
			Skip:  c.flags.Type != "",
		},
		prompt.StringPrompt{
			Label:   "Cloud provider",
			Default: AwsProvider,
			Validate: func(input string) error {
				switch input {
				case AwsProvider:
					return nil
				default:
					return errors.New("unsupported cloud provider")
				}
			},
			Value: &c.flags.Provider,
			Skip:  c.flags.Provider != "",
		},
		prompt.StringPrompt{
			Label:   "Region",
			Default: "us-east-1",
			Value:   &c.flags.Region,
			Skip:    c.flags.Region != "",
		},
		prompt.ConditionalPrompt{
			If: prompt.BoolPrompt{
				Label: "Does your environment require configuration",
			},
			Then: prompt.MapPrompt{
				KeyPrompt: prompt.StringPrompt{
					Label: "\tConfiguration key",
				},
				ValuePrompt: prompt.StringPrompt{
					Label: "\tConfiguration value",
				},
				Value:     c.envCfg,
				NextLabel: "Add another configuration",
			},
			Skip: c.flags.Config != "",
		},
	})
	if err != nil {
		return err
	}

	c.showEventConfirmation()

	var create bool
	_, err = prompt.BoolPrompt{
		Label: "Create this environment",
		Value: &create,
	}.Show(ctx)

	if err != nil {
		return err
	}

	if !create {
		return builder.NewErrPromptAbort(
			errors.New("\nTo view all options for creating a Meroxa Environment,\n " +
				"please run \"meroxa help env create\""),
		)
	}
	return nil
}

func (c *Create) Usage() string {
	return "create NAME"
}

func (c *Create) Docs() builder.Docs {
	return builder.Docs{
		Short: "Create an environment",
		Example: `
meroxa env create my-env --type hosted --provider aws --region us-east-1 --config aws_access_key_id=1234 --config aws_secret_access_key=5678
`,
	}
}
