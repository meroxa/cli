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
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs             = (*Remove)(nil)
	_ builder.CommandWithAliases          = (*Remove)(nil)
	_ builder.CommandWithArgs             = (*Remove)(nil)
	_ builder.CommandWithClient           = (*Remove)(nil)
	_ builder.CommandWithLogger           = (*Remove)(nil)
	_ builder.CommandWithExecute          = (*Remove)(nil)
	_ builder.CommandWithConfirmWithValue = (*Remove)(nil)
)

type removeEnvironmentClient interface {
	DeleteEnvironment(ctx context.Context, nameOrUUID string) (*meroxa.Environment, error)
}

type Remove struct {
	client removeEnvironmentClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}
}

func (r *Remove) Usage() string {
	return "remove NAMEorUUID"
}

func (r *Remove) Docs() builder.Docs {
	return builder.Docs{
		Short: "Remove environment",
	}
}

func (r *Remove) ValueToConfirm(_ context.Context) (wantInput string) {
	return r.args.NameOrUUID
}

func (r *Remove) Execute(ctx context.Context) error {
	r.logger.Infof(ctx, "Environment %q is being removed...", r.args.NameOrUUID)

	e, err := r.client.DeleteEnvironment(ctx, r.args.NameOrUUID)
	if err != nil {
		return err
	}

	r.logger.Infof(ctx, "Run `meroxa env describe %s` for status.", r.args.NameOrUUID)
	r.logger.JSON(ctx, e)

	return nil
}

func (r *Remove) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Remove) Client(client meroxa.Client) {
	r.client = client
}

func (r *Remove) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires environment name")
	}

	r.args.NameOrUUID = args[0]
	return nil
}

func (r *Remove) Aliases() []string {
	return []string{"rm", "delete"}
}
