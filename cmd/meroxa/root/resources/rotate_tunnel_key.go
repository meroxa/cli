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
	_ builder.CommandWithDocs             = (*RotateTunnelKey)(nil)
	_ builder.CommandWithArgs             = (*RotateTunnelKey)(nil)
	_ builder.CommandWithFlags            = (*RotateTunnelKey)(nil)
	_ builder.CommandWithClient           = (*RotateTunnelKey)(nil)
	_ builder.CommandWithLogger           = (*RotateTunnelKey)(nil)
	_ builder.CommandWithExecute          = (*RotateTunnelKey)(nil)
	_ builder.CommandWithConfirmWithValue = (*RotateTunnelKey)(nil)
)

type rotateKeyActionClient interface {
	RotateTunnelKeyForResource(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
	GetResourceByName(ctx context.Context, name string) (*meroxa.Resource, error)
}

type RotateTunnelKey struct {
	client rotateKeyActionClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
	}
}

func (u *RotateTunnelKey) Usage() string {
	return "rotate-tunnel-key NAME"
}

func (u *RotateTunnelKey) Docs() builder.Docs {
	return builder.Docs{
		Short: "Rotate the tunnel key for a resource",
		Long:  "Rotate the tunnel key for a Meroxa resource.",
	}
}

func (u *RotateTunnelKey) ValueToConfirm(ctx context.Context) string {
	u.logger.Infof(ctx, "Rotating tunnel key will restart the tunnel and disconnect existing connections.")
	return u.args.Name
}

func (u *RotateTunnelKey) Execute(ctx context.Context) error {
	r, err := u.client.RotateTunnelKeyForResource(ctx, u.args.Name)
	if err != nil {
		return err
	}

	u.logger.Infof(ctx, "Resource %q tunnel key is successfully rotated!", u.args.Name)
	if tun := r.SSHTunnel; tun != nil {
		u.logger.Info(ctx, "Paste the following public key on your host:")
		u.logger.Info(ctx, tun.PublicKey)
		u.logger.Info(ctx, "Meroxa will try to connect to the resource for 60 minutes and send an email confirmation after a successful resource validation.") //nolint
	}
	u.logger.JSON(ctx, r)

	return nil
}

func (u *RotateTunnelKey) Flags() []builder.Flag {
	return builder.BuildFlags(&u.flags)
}

func (u *RotateTunnelKey) Logger(logger log.Logger) {
	u.logger = logger
}

func (u *RotateTunnelKey) Client(client meroxa.Client) {
	u.client = client
}

func (u *RotateTunnelKey) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires resource name")
	}

	u.args.Name = args[0]
	return nil
}
