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

	"github.com/meroxa/cli/cmd/meroxa/builder"
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
)

const (
	DefaultProvider = "aws"
	DefaultType     = "dedicated"
	DefaultRegion   = "us-east-1"
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
		Type          string `long:"type" usage:"environment type (\"dedicated\", when not specified. It can also be \"hosted\")"`
		Provider      string `long:"provider" usage:"environment cloud provider to use (e.g.: AWS)"`
		Region        string `long:"region" usage:"environment region (e.g.: us-east-1)"`
		Configuration string `long:"config" usage:"environment configuration based on type and provider (e.g.: --config aws_access_key_id=my_access_key)"`
		Interactive   bool   `long:"interactive" short:"i" usage:"Interactive mode""`
	}
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

func (c *Create) userSpecifiedValues() bool {
	return c.flags.Type == "" || c.flags.Provider == "" || c.flags.Region == "" || c.flags.Configuration == ""
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
}

func (c *Create) Execute(ctx context.Context) error {
	c.logger.Infof(ctx, "%b", c.flags.Interactive)
	if c.flags.Interactive {
		// TODO: Call interactive mode and set values

	}

	//e := &meroxa.CreateEnvironmentInput{}
	//c.setUserValues(e)
	//
	//// Show specified values (either default or specified by user)
	//c.logger.Infof(ctx, "We are going to create an environment that will look like this:\n"+
	//	"\t type: %q\n"+
	//	"\t provider: %q\n"+
	//	"\t region: %q", DefaultType, DefaultProvider, DefaultRegion)
	//
	//prompt := promptui.Prompt{
	//	Label:     "Do you want to proceed?",
	//	IsConfirm: true,
	//}
	//
	//_, error := prompt.Run()
	//
	//if error != nil {
	//	c.logger.Infof(ctx, "If you want to configure an environment with different settings,\n "+
	//		"please run \"meroxa help env create\". \n"+
	//		"For a more guided approach, run in interactive mode: \"meroxa env create -i\"")
	//}
	//
	//if c.flags.Configuration != "" {
	//	var config map[string]interface{}
	//	err := json.Unmarshal([]byte(c.flags.Configuration), &config)
	//	if err != nil {
	//		return fmt.Errorf("could not parse configuration: %w", err)
	//	}
	//
	//	e.Configuration = config
	//}
	//
	//c.logger.Infof(ctx, "Provisioning environment...")
	//
	//environment, err := c.client.CreateEnvironment(ctx, e)
	//
	//if err != nil {
	//	return err
	//}
	//
	//c.logger.Infof(ctx, "Environment %q is being provisioned. Run `meroxa env describe %q` for status", environment.Name, environment.Name)
	//c.logger.JSON(ctx, environment)

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
