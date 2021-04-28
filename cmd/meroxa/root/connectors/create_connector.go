package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
)

var (
	_ builder.CommandWithDocs    = (*CreateConnector)(nil)
	_ builder.CommandWithArgs    = (*CreateConnector)(nil)
	_ builder.CommandWithFlags   = (*CreateConnector)(nil)
	_ builder.CommandWithClient  = (*CreateConnector)(nil)
	_ builder.CommandWithLogger  = (*CreateConnector)(nil)
	_ builder.CommandWithExecute = (*CreateConnector)(nil)
)

type createConnectorClient interface {
	GetResourceByName(ctx context.Context, name string) (*meroxa.Resource, error)
	CreateConnector(ctx context.Context, input meroxa.CreateConnectorInput) (*meroxa.Connector, error)
}

type CreateConnector struct {
	client createConnectorClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Input       string `long:"input"      short:""  usage:"command delimited list of input streams" required:"true"`
		Config      string `long:"config"      short:"c"  usage:"connector configuration"`
		Metadata    string `long:"metadata"    short:"m" usage:"connector metadata" hidden:"true"`
		Source      string `long:"from"    short:"" usage:"resource name to use as source"`
		Destination string `long:"to"    short:"" usage:"resource name to use as destination"`
		Pipeline    string `long:"pipeline"    short:"" usage:"pipeline name to attach connector to"`
	}
}

func (c *CreateConnector) parseJSONMap(str string) (out map[string]interface{}, err error) {
	out = make(map[string]interface{})
	if str != "" {
		err = json.Unmarshal([]byte(str), &out)
	}
	return out, err
}

func (c *CreateConnector) Execute(ctx context.Context) error {
	// TODO: Implement something like dependant flags in Builder
	if c.flags.Source == "" && c.flags.Destination == "" {
		return errors.New("requires either a source (--from) or a destination (--to)")
	}

	config, err := c.parseJSONMap(c.flags.Config)
	if err != nil {
		return errors.New("can't parse config, make sure it is a valid JSON map")
	}

	metadata, err := c.parseJSONMap(c.flags.Metadata)
	if err != nil {
		return errors.New("can't parse metadata, make sure it is a valid JSON map")
	}

	// merge in input
	config["input"] = c.flags.Input

	// merge in connector type
	var resourceName string
	switch {
	case c.flags.Source != "":
		resourceName = c.flags.Source
		metadata["mx:connectorType"] = "source"
	case c.flags.Destination != "":
		resourceName = c.flags.Destination
		metadata["mx:connectorType"] = "destination"
	default:
		return errors.New("requires either a source (--from) or a destination (--to)")
	}

	res, err := c.client.GetResourceByName(ctx, resourceName)
	if err != nil {
		return fmt.Errorf("can't fetch resource with name %q: %w", resourceName, err)
	}

	switch {
	case c.flags.Source != "":
		c.logger.Infof(ctx, "Creating connector from source %s...\n", resourceName)
	case c.flags.Destination != "":
		c.logger.Infof(ctx, "Creating connector to destination %s...\n", resourceName)
	}

	connector, err := c.client.CreateConnector(ctx, meroxa.CreateConnectorInput{
		Name:          c.args.Name,
		ResourceID:    res.ID,
		PipelineName:  c.flags.Pipeline,
		Configuration: config,
		Metadata:      metadata,
	})

	if err != nil {
		return err
	}

	c.logger.Infof(ctx, "Connector %s successfully created!\n", connector.Name)
	c.logger.JSON(ctx, connector)

	return nil
}

func (c *CreateConnector) Client(client *meroxa.Client) {
	c.client = client
}

func (c *CreateConnector) Logger(logger log.Logger) {
	c.logger = logger
}

func (c *CreateConnector) Flags() []builder.Flag {
	return builder.BuildFlags(&c.flags)
}

func (c *CreateConnector) ParseArgs(args []string) error {
	if len(args) > 0 {
		c.args.Name = args[0]
	}
	return nil
}

func (c *CreateConnector) Usage() string {
	return "create [NAME] [flags]"
}

func (c *CreateConnector) Docs() builder.Docs {
	return builder.Docs{
		Short: "Create a connector",
		Long:  "Use `connectors create` to create a connector from a source (--from) or to a destination (--to)",
		Example: "\n" +
			"meroxa connectors create [NAME] --from pg2kafka --input accounts \n" +
			"meroxa connectors create [NAME] --to pg2redshift --input orders # --input will be the desired stream \n" +
			"meroxa connectors create [NAME] --to pg2redshift --input orders --pipeline my-pipeline\n",
	}
}
