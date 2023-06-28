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
	"github.com/davecgh/go-spew/spew"
	"github.com/meroxa/turbine-core/pkg/ir"
	"path/filepath"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/flink"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
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
		Jar     string   `long:"jar" required:"true" usage:"Path to Flink Job jar file"`
		Secrets []string `short:"s" long:"secret" usage:"environment variables to inject into the Flink Job (e.g.: --secret API_KEY=$API_KEY --secret ACCESS_KEY=abcdef)"` //nolint:lll
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

	if filepath.Ext(jarPath) != ".jar" {
		return fmt.Errorf("please provide a JAR file to the --jar flag")
	}

	secrets := utils.StringSliceToStringMap(d.flags.Secrets)
	spec, err := flink.GetIRSpec(ctx, jarPath, secrets, d.logger)
	if err != nil {
		fmt.Printf("failed to extract IR spec: %v\n", err)
		// non-blocking as of yet... is this still true?
	}
	spec.Definition.Metadata.SpecVersion = ir.LatestSpecVersion // temporary workaround

	fmt.Printf("spec:\n%s\n", spew.Sdump(spec))

	name := d.args.Name
	if name == "" {
		return fmt.Errorf("the name of your Flink Job be provided as an argument")
	}

	filename := filepath.Base(jarPath)
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
	input := &meroxa.CreateFlinkJobInput{Name: name, JarURL: source.GetUrl}
	if spec != nil {
		d.logger.StartSpinner("\t", "Adding Meroxa integrations to request...")
		bytes, err := spec.Marshal()
		//bytes, err := json.Marshal(spec)
		if err != nil {
			d.logger.Errorf(ctx, "\t 𐄂 Unable to add Meroxa integrations to request")
			d.logger.StopSpinnerWithStatus("\t", log.Failed)
			return err
		}
		input.Spec = string(bytes)
		input.SpecVersion = spec.Definition.Metadata.SpecVersion
	}
	fmt.Printf("GetUrl: %s\n", source.GetUrl)
	fj, err := d.client.CreateFlinkJob(ctx, input)
	if err != nil {
		d.logger.Errorf(ctx, "\t 𐄂 Unable to create Flink job")
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	d.logger.StopSpinnerWithStatus("Flink job created", log.Successful)
	d.logger.JSON(ctx, fj)
	return nil
}
