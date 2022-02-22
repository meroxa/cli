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

package apps

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
)

const (
	dockerHubUserNameEnv    = "DOCKER_HUB_USERNAME"
	dockerHubAccessTokenEnv = "DOCKER_HUB_ACCESS_TOKEN" // nolint:gosec
)

type Deploy struct {
	flags struct {
		Path string `long:"path" description:"path to the app directory"`
	}

	config   config.Config
	logger   log.Logger
	path     string
	goDeploy turbineCLI.GoDeploy
}

var (
	_ builder.CommandWithDocs    = (*Deploy)(nil)
	_ builder.CommandWithFlags   = (*Deploy)(nil)
	_ builder.CommandWithExecute = (*Deploy)(nil)
	_ builder.CommandWithLogger  = (*Deploy)(nil)
	_ builder.CommandWithConfig  = (*Deploy)(nil)
)

func (*Deploy) Usage() string {
	return "deploy"
}

func (*Deploy) Docs() builder.Docs {
	return builder.Docs{
		Short: "Deploy a Meroxa Data Application",
		Example: "meroxa apps deploy # assumes you run it from the app directory\n" +
			"meroxa apps deploy --path ./my-app",
	}
}

func (d *Deploy) Config(cfg config.Config) {
	d.config = cfg
}

func (d *Deploy) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Deploy) getDockerHubUserNameEnv() string {
	if v := os.Getenv(dockerHubUserNameEnv); v != "" {
		return v
	}
	return d.config.GetString(dockerHubUserNameEnv)
}

func (d *Deploy) getDockerHubAccessTokenEnv() string {
	if v := os.Getenv(dockerHubAccessTokenEnv); v != "" {
		return v
	}
	return d.config.GetString(dockerHubAccessTokenEnv)
}

func (d *Deploy) checkRequiredEnvVars() error {
	if d.getDockerHubUserNameEnv() == "" || d.getDockerHubAccessTokenEnv() == "" {
		msg := fmt.Sprintf("both %q and %q are required to be set to deploy your application", dockerHubUserNameEnv, dockerHubAccessTokenEnv)
		return errors.New(msg)
	}

	d.goDeploy.DockerHubUserNameEnv = d.getDockerHubUserNameEnv()
	d.goDeploy.DockerHubAccessTokenEnv = d.getDockerHubAccessTokenEnv()

	return nil
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Deploy) Execute(ctx context.Context) error {
	err := d.checkRequiredEnvVars()
	if err != nil {
		return err
	}

	d.path = turbineCLI.GetPath(d.flags.Path)
	lang, err := turbineCLI.GetLangFromAppJSON(d.path)
	if err != nil {
		return err
	}

	switch lang {
	case "go", GoLang:
		return d.goDeploy.DeployGoApp(ctx, d.path, d.logger)
	case "js", JavaScript, NodeJs:
		return turbineCLI.DeployJSApp(ctx, d.path, d.logger)
	default:
		return fmt.Errorf("language %q not supported. Currently, we support \"javascript\" and \"go\"", lang)
	}
}
