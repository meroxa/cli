/*
Copyright © 2022 Meroxa Inc

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
	"fmt"
	"path/filepath"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

type deployFlinkJobClient interface {
	CreateSourceV2(ctx context.Context, input *meroxa.CreateSourceInputV2) (*meroxa.Source, error)
	CreateFlinkJob(ctx context.Context, input *meroxa.CreateFlinkJobInput) (*meroxa.FlinkJob, error)
}

type Deploy struct {
	args struct {
		Name string
	}

	flags struct {
		Jar string `long:"jar" required:"true" usage:"Path to Flink Job jar file"`
	}

	client deployFlinkJobClient
	config config.Config
	logger log.Logger
}

var (
	_ builder.CommandWithClient  = (*Deploy)(nil)
	_ builder.CommandWithConfig  = (*Deploy)(nil)
	_ builder.CommandWithDocs    = (*Deploy)(nil)
	_ builder.CommandWithExecute = (*Deploy)(nil)
	_ builder.CommandWithArgs    = (*Deploy)(nil)
	_ builder.CommandWithFlags   = (*Deploy)(nil)
	_ builder.CommandWithLogger  = (*Deploy)(nil)
)

func (*Deploy) Usage() string {
	return "deploy NAME --jar /home/job.jar"
}

func (*Deploy) Docs() builder.Docs {
	return builder.Docs{
		Short: "Deploy a Flink Job",
	}
}

func (d *Deploy) Config(cfg config.Config) {
	d.config = cfg
}

func (d *Deploy) Client(client meroxa.Client) {
	d.client = client
}

func (d *Deploy) ParseArgs(args []string) error {
	if len(args) > 0 {
		d.args.Name = args[0]
	}
	return nil
}

func (d *Deploy) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Deploy) Execute(ctx context.Context) error {
	jarPath := d.flags.Jar
	if jarPath == "" {
		return fmt.Errorf("the path to your Flink Job jar file must be provided to the --jar flag")
	}

	filename := filepath.Base(jarPath)

	name := d.args.Name
	if name == "" {
		return fmt.Errorf("the name of your Flink Job be provided as an argument")
	}

	d.logger.StartSpinner("\t", "Fetching Meroxa Platform source...")
	source, err := d.client.CreateSourceV2(ctx, &meroxa.CreateSourceInputV2{Filename: filename})
	if err != nil {
		d.logger.Errorf(ctx, "\t 𐄂 Unable to fetch source")
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}
	d.logger.StopSpinnerWithStatus("Platform source fetched", log.Successful)

	err = turbine.UploadFile(ctx, d.logger, jarPath, source.PutUrl)
	if err != nil {
		return err
	}

	d.logger.StartSpinner("\t", "Creating Flink job...")
	fj, err := d.client.CreateFlinkJob(ctx, &meroxa.CreateFlinkJobInput{Name: name, JarURL: source.GetUrl})
	if err != nil {
		d.logger.Errorf(ctx, "\t 𐄂 Unable to create Flink job")
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	d.logger.StopSpinnerWithStatus("Flink job created", log.Successful)
	d.logger.JSON(ctx, fj)
	return nil
}