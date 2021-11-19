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

package resources

import (
	"context"
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

type removeResourceClient interface {
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
	DeleteResource(ctx context.Context, nameOrID string) error
}

type Remove struct {
	client removeResourceClient
	logger log.Logger

	args struct {
		NameOrID string
	}
}

func (r *Remove) Usage() string {
	return "remove NAME"
}

func (r *Remove) Docs() builder.Docs {
	return builder.Docs{
		Short: "Remove resource",
	}
}

func (r *Remove) ValueToConfirm(_ context.Context) (wantInput string) {
	return r.args.NameOrID
}

func (r *Remove) Execute(ctx context.Context) error {
	r.logger.Infof(ctx, "Removing resource %q...", r.args.NameOrID)

	res, err := r.client.GetResourceByNameOrID(ctx, r.args.NameOrID)
	if err != nil {
		return err
	}

	err = r.client.DeleteResource(ctx, r.args.NameOrID)

	if err != nil {
		return err
	}

	r.logger.Infof(ctx, "Resource %q successfully removed", r.args.NameOrID)
	r.logger.JSON(ctx, res)

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
		return errors.New("requires resource name")
	}

	r.args.NameOrID = args[0]
	return nil
}

func (r *Remove) Aliases() []string {
	return []string{"rm", "delete"}
}

var (
	_ builder.CommandWithDocs             = (*Remove)(nil)
	_ builder.CommandWithAliases          = (*Remove)(nil)
	_ builder.CommandWithArgs             = (*Remove)(nil)
	_ builder.CommandWithClient           = (*Remove)(nil)
	_ builder.CommandWithLogger           = (*Remove)(nil)
	_ builder.CommandWithExecute          = (*Remove)(nil)
	_ builder.CommandWithConfirmWithValue = (*Remove)(nil)
)
