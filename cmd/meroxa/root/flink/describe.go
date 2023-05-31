/*
Copyright © 2023 Meroxa Inc

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

package flink

import (
	"context"
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Describe)(nil)
	_ builder.CommandWithArgs    = (*Describe)(nil)
	_ builder.CommandWithClient  = (*Describe)(nil)
	_ builder.CommandWithLogger  = (*Describe)(nil)
	_ builder.CommandWithExecute = (*Describe)(nil)
)

type getJobClient interface {
	GetFlinkJob(ctx context.Context, nameOrUUID string) (*meroxa.FlinkJob, error)
}

type Describe struct {
	client getJobClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}
}

func (d *Describe) Usage() string {
	return "describe"
}

func (d *Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe the details of a Flink Job",
	}
}

func (d *Describe) Execute(ctx context.Context) error {
	flinkJob, err := d.client.GetFlinkJob(ctx, d.args.NameOrUUID)
	if err != nil {
		return err
	}

	d.logger.JSON(ctx, flinkJob)
	d.logger.Info(ctx, display.FlinkJobTable(flinkJob))

	return nil
}

func (d *Describe) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Describe) Client(client meroxa.Client) {
	d.client = client
}

func (d *Describe) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires Flink Job name or UUID")
	}

	d.args.NameOrUUID = args[0]
	return nil
}
