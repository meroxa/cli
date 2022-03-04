package functions

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/volatiletech/null/v8"

	"github.com/meroxa/cli/cmd/meroxa/builder"
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

type createFunctionClient interface {
	CreateFunction(ctx context.Context, input *meroxa.CreateFunctionInput) (*meroxa.Function, error)
}

type Create struct {
	client createFunctionClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		InputStream string   `long:"input-stream" usage:"an input stream to the function" required:"true"`
		Image       string   `long:"image" usage:"Docker image name" required:"true"`
		Command     string   `long:"command" usage:"Entrypoint command"`
		Args        string   `long:"args" usage:"Arguments to the entrypoint"`
		EnvVars     []string `long:"env" usage:"List of environment variables to set in the function"`
		Pipeline    string   `long:"pipeline" usage:"pipeline name to attach function to" required:"true"`
		Application string   `long:"app" usage:"application name or UUID to which this function belongs" required:"true"`
	}
}

func (c *Create) Usage() string {
	return "create [NAME] [flags]"
}

func (c *Create) Docs() builder.Docs {
	return builder.Docs{
		Short: "Create a function",
		Long:  "Use `functions create` to create a function to process records from an input stream (--input-stream)",
		Example: `
meroxa functions create [NAME] --input-stream connector-output-stream --image myimage --app my-app
meroxa functions create [NAME] --input-stream connector-output-stream --image myimage --app my-app --env FOO=BAR --env BAR=BAZ
`,
	}
}

func (c *Create) Execute(ctx context.Context) error {
	envVars, err := c.parseEnvVars(c.flags.EnvVars)
	if err != nil {
		return err
	}

	var (
		command []string
		args    []string
	)
	if cmd := c.flags.Command; cmd != "" {
		command, err = shellwords.Parse(cmd)
		if err != nil {
			return err
		}
	}
	if a := c.flags.Args; a != "" {
		args, err = shellwords.Parse(a)
		if err != nil {
			return err
		}
	}

	fun, err := c.client.CreateFunction(
		ctx,
		&meroxa.CreateFunctionInput{
			Name:        c.args.Name,
			InputStream: c.flags.InputStream,
			Pipeline: meroxa.PipelineIdentifier{
				Name: null.StringFrom(c.flags.Pipeline),
			},
			Application: meroxa.ApplicationIdentifier{
				Name: null.StringFrom(c.flags.Application),
			},
			Image:   c.flags.Image,
			Command: command,
			Args:    args,
			EnvVars: envVars,
		},
	)
	if err != nil {
		return err
	}

	c.logger.Infof(ctx, "Function %q successfully created!\n", fun.Name)
	c.logger.JSON(ctx, fun)

	return nil
}

func (c *Create) parseEnvVars(envVars []string) (map[string]string, error) {
	m := make(map[string]string)
	for _, ev := range envVars {
		var (
			split = strings.SplitN(ev, "=", 2) //nolint
			key   = split[0]
			val   = split[1]
		)

		if key == "" || val == "" {
			return nil, fmt.Errorf("error parsing env var %q", ev)
		}

		m[key] = val
	}

	return m, nil
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
