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
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

var (
	_ builder.CommandWithDocs    = (*Describe)(nil)
	_ builder.CommandWithArgs    = (*Describe)(nil)
	_ builder.CommandWithClient  = (*Describe)(nil)
	_ builder.CommandWithLogger  = (*Describe)(nil)
	_ builder.CommandWithExecute = (*Describe)(nil)
)

type describeResourceClient interface {
	GetResourceByName(ctx context.Context, name string) (*meroxa.Resource, error)
}

type Describe struct {
	client describeResourceClient
	logger log.Logger

	args struct {
		Name string
	}
}

func (d *Describe) Usage() string {
	return "describe [NAME]"
}

func (d *Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe resource",
	}
}

func (d *Describe) Execute(ctx context.Context) error {
	resource, err := d.client.GetResourceByName(ctx, d.args.Name)
	if err != nil {
		return err
	}

	d.logger.Info(ctx, utils.ResourcesTable([]*meroxa.Resource{resource}))
	d.logger.JSON(ctx, resource)

	return nil
}

func (d *Describe) Client(client *meroxa.Client) {
	d.client = client
}

func (d *Describe) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Describe) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires resource name")
	}

	d.args.Name = args[0]
	return nil
}
