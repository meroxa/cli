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

package pipelines

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/root/environments"
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

type createPipelineClient interface {
	CreatePipeline(ctx context.Context, input *meroxa.CreatePipelineInput) (*meroxa.Pipeline, error)
}

type Create struct {
	client createPipelineClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Metadata string `long:"metadata" short:"m" usage:"pipeline metadata"`
		// TODO: Add support to builder to create flags with an alias (--env | --environment)
		Environment string `long:"env" usage:"environment (name or UUID) where pipeline will be created"`
	}
}

func (c *Create) Execute(ctx context.Context) error {
	var env string

	p := &meroxa.CreatePipelineInput{
		Name: c.args.Name,
	}

	if c.flags.Metadata != "" {
		var metadata map[string]interface{}
		err := json.Unmarshal([]byte(c.flags.Metadata), &metadata)
		if err != nil {
			return fmt.Errorf("could not parse metadata: %w", err)
		}

		p.Metadata = metadata
	}

	// If the environment specified is not the common environment.
	if c.flags.Environment != "" && c.flags.Environment != string(meroxa.EnvironmentTypeCommon) {
		err := builder.CheckFeatureFlag(c, &environments.Environments{})
		if err != nil {
			return err
		}

		env = c.flags.Environment
		p.Environment = &meroxa.EnvironmentIdentifier{}

		_, err = uuid.Parse(c.flags.Environment)

		if err == nil {
			p.Environment.UUID = c.flags.Environment
		} else {
			p.Environment.Name = c.flags.Environment
		}
	} else {
		env = string(meroxa.EnvironmentTypeCommon)
	}

	c.logger.Infof(ctx, "Creating pipeline %q in %q environment...", c.args.Name, env)
	pipeline, err := c.client.CreatePipeline(ctx, p)

	if err != nil {
		return err
	}

	c.logger.Infof(ctx, "Pipeline %q successfully created!", c.args.Name)
	c.logger.JSON(ctx, pipeline)

	return nil
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
		return errors.New("requires a pipeline name")
	}

	c.args.Name = args[0]
	return nil
}

func (c *Create) Usage() string {
	return "create NAME"
}

func (c *Create) Docs() builder.Docs {
	return builder.Docs{
		Short: "Create a pipeline",
	}
}
