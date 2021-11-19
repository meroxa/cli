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

var (
	_ builder.CommandWithDocs    = (*Validate)(nil)
	_ builder.CommandWithArgs    = (*Validate)(nil)
	_ builder.CommandWithFlags   = (*Validate)(nil)
	_ builder.CommandWithClient  = (*Validate)(nil)
	_ builder.CommandWithLogger  = (*Validate)(nil)
	_ builder.CommandWithExecute = (*Validate)(nil)
)

type validateResourceClient interface {
	ValidateResource(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
}

type Validate struct {
	client validateResourceClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
	}
}

func (u *Validate) Usage() string {
	return "validate NAME"
}

func (u *Validate) Docs() builder.Docs {
	return builder.Docs{
		Short: "Validate a resource",
		Long:  "Validate a Meroxa resource.",
	}
}

func (u *Validate) Execute(ctx context.Context) error {
	r, err := u.client.ValidateResource(ctx, u.args.Name)
	if err != nil {
		return err
	}

	if r.SSHTunnel == nil {
		u.logger.Infof(ctx, "Resource %q is successfully validated!", u.args.Name)
	} else {
		u.logger.Infof(ctx, "Resource %q validation is successfully triggered!", u.args.Name)
		u.logger.Info(ctx, "Meroxa will try to connect to the resource for 60 minutes and send an email confirmation after a successful resource validation.") //nolint
	}

	u.logger.JSON(ctx, r)
	return nil
}

func (u *Validate) Flags() []builder.Flag {
	return builder.BuildFlags(&u.flags)
}

func (u *Validate) Logger(logger log.Logger) {
	u.logger = logger
}

func (u *Validate) Client(client meroxa.Client) {
	u.client = client
}

func (u *Validate) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires resource name")
	}

	u.args.Name = args[0]
	return nil
}
