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
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/flink"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/turbine-core/pkg/ir"
)

type Deploy struct {
	args struct {
		Name string
	}

	flags struct {
		Jar     string   `long:"jar" required:"true" usage:"Path to Flink Job jar file"`
		Secrets []string `short:"s" long:"secret" usage:"environment variables to inject into the Flink Job (e.g.: --secret API_KEY=$API_KEY --secret ACCESS_KEY=abcdef)"` //nolint:lll
	}

	client global.BasicClient
	config config.Config
	logger log.Logger
}

var (
	_ builder.CommandWithBasicClient = (*Deploy)(nil)
	_ builder.CommandWithConfig      = (*Deploy)(nil)
	_ builder.CommandWithDocs        = (*Deploy)(nil)
	_ builder.CommandWithExecute     = (*Deploy)(nil)
	_ builder.CommandWithArgs        = (*Deploy)(nil)
	_ builder.CommandWithFlags       = (*Deploy)(nil)
	_ builder.CommandWithLogger      = (*Deploy)(nil)
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

func (d *Deploy) BasicClient(client global.BasicClient) {
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

	name := d.args.Name
	if name == "" {
		return fmt.Errorf("the name of your Flink Job be provided as an argument")
	}

	secrets := utils.StringSliceToStringMap(d.flags.Secrets)
	_, err := flink.GetIRSpec(ctx, jarPath, secrets, d.logger)
	if err != nil {
		d.logger.Warnf(ctx, "failed to extract IR spec: %v\n", err)
		// non-blocking
	}

	d.logger.StartSpinner("\t", "Fetching Meroxa Platform source...")

	d.logger.StopSpinnerWithStatus("Platform source fetched", log.Successful)

	// Logging happens inside UploadFile

	// Creqte flink job

	d.logger.StopSpinnerWithStatus("Flink job created", log.Successful)
	return nil
}

func (d *Deploy) addIntegrations(ctx context.Context, spec *ir.DeploymentSpec) error {
	d.logger.StartSpinner("\t", "Checking Meroxa integrations...")
	successMsg := "Finished checking Meroxa integrations"
	if spec != nil {
		var bytes []byte
		bytes, err := json.Marshal(spec)
		if err != nil {
			d.logger.Errorf(ctx, "\t 𐄂 Unable to add Meroxa integrations to request")
			d.logger.StopSpinnerWithStatus("\t", log.Failed)
			return err
		}
		var inputSpec map[string]interface{}
		if unmarshalErr := json.Unmarshal(bytes, &inputSpec); unmarshalErr != nil {
			d.logger.Errorf(ctx, "\t 𐄂 Unable to add Meroxa integrations to request")
			d.logger.StopSpinnerWithStatus("\t", log.Failed)
			return unmarshalErr
		}
		successMsg = "Added Meroxa integrations to request"

	}
	d.logger.StopSpinnerWithStatus(successMsg, log.Successful)
	return nil
}
