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
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
)

var (
	_ builder.CommandWithDocs    = (*Remove)(nil)
	_ builder.CommandWithAliases = (*Remove)(nil)
	_ builder.CommandWithArgs    = (*Remove)(nil)
	_ builder.CommandWithClient  = (*Remove)(nil)
	_ builder.CommandWithLogger  = (*Remove)(nil)
	_ builder.CommandWithExecute = (*Remove)(nil)
	_ builder.CommandWithConfirm = (*Remove)(nil)
)

type removeEnvironmentClient interface {
	DeleteEnvironment(ctx context.Context, nameOrUUID string) error
	// GetEnvironment TODO: Remove this one if Env can be returned by Delete
	GetEnvironment(ctx context.Context, nameOrUUID string) (*meroxa.Environment, error)
}

type Remove struct {
	client removeEnvironmentClient
	logger log.Logger

	args struct {
		Name string
	}
}

func (r *Remove) Usage() string {
	return "remove NAME"
}

func (r *Remove) Docs() builder.Docs {
	return builder.Docs{
		Short: "Remove environment",
	}
}

func (r *Remove) Confirm(_ context.Context) (wantInput string) {
	return r.args.Name
}

func (r *Remove) Execute(ctx context.Context) error {
	r.logger.Infof(ctx, "Removing environment %q...", r.args.Name)

	// TODO: Check if this could return also env
	err := r.client.DeleteEnvironment(ctx, r.args.Name)
	if err != nil {
		return err
	}

	r.logger.Infof(ctx, "Environment %q successfully removed", r.args.Name)

	e, err := r.client.GetEnvironment(ctx, r.args.Name)
	if err != nil {
		return err
	}
	r.logger.JSON(ctx, e)

	return nil
}

func (r *Remove) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Remove) Client(client *meroxa.Client) {
	r.client = client
}

func (r *Remove) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires environment name")
	}

	r.args.Name = args[0]
	return nil
}

func (r *Remove) Aliases() []string {
	return []string{"rm", "delete"}
}
