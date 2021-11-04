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
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

type updatePipelineClient interface {
	GetPipelineByName(ctx context.Context, name string) (*meroxa.Pipeline, error)
	UpdatePipelineStatus(ctx context.Context, pipelineID int, state meroxa.Action) (*meroxa.Pipeline, error)
	UpdatePipeline(ctx context.Context, pipelineID int, pipeline *meroxa.UpdatePipelineInput) (*meroxa.Pipeline, error)
}

type Update struct {
	client updatePipelineClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		State string `long:"state" usage:"new pipeline state (pause | resume | restart)"`
		Name  string `long:"name" usage:"new pipeline name"`
	}
}

func (u *Update) Usage() string {
	return "update NAME"
}

func (u *Update) Docs() builder.Docs {
	return builder.Docs{
		Short: "Update pipeline name, state or metadata",
		Example: "\n" +
			"meroxa pipeline update old-name --name new-name\n" +
			"meroxa pipeline update pipeline-name --state pause\n" +
			"meroxa pipeline update pipeline-name --metadata '{\"key\":\"value\"}'\n" +
			"meroxa pipeline update pipeline-name --state restart",
	}
}

func (u *Update) Execute(ctx context.Context) error {
	// TODO: Implement something like dependent flags in Builder
	if u.flags.Name == "" && u.flags.State == "" {
		return errors.New("requires either --name, --state or --metadata")
	}

	u.logger.Infof(ctx, "Updating pipeline %q...", u.args.Name)

	p, err := u.client.GetPipelineByName(ctx, u.args.Name)

	if err != nil {
		return err
	}

	// update state/status separately
	if u.flags.State != "" {
		p, err = u.client.UpdatePipelineStatus(ctx, p.ID, meroxa.Action(u.flags.State))
		if err != nil {
			return err
		}
	}

	// call meroxa-go to update either name or metadata
	if u.flags.Name != "" {
		pi := &meroxa.UpdatePipelineInput{
			Name: u.flags.Name,
		}

		p, err = u.client.UpdatePipeline(ctx, p.ID, pi)
		if err != nil {
			return err
		}
	}

	u.logger.Infof(ctx, "Pipeline %q successfully updated!", u.args.Name)
	u.logger.JSON(ctx, p)
	return nil
}

func (u *Update) Flags() []builder.Flag {
	return builder.BuildFlags(&u.flags)
}

func (u *Update) Logger(logger log.Logger) {
	u.logger = logger
}

func (u *Update) Client(client meroxa.Client) {
	u.client = client
}

func (u *Update) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires pipeline name")
	}

	u.args.Name = args[0]
	return nil
}

var (
	_ builder.CommandWithDocs    = (*Update)(nil)
	_ builder.CommandWithArgs    = (*Update)(nil)
	_ builder.CommandWithFlags   = (*Update)(nil)
	_ builder.CommandWithClient  = (*Update)(nil)
	_ builder.CommandWithLogger  = (*Update)(nil)
	_ builder.CommandWithExecute = (*Update)(nil)
)
