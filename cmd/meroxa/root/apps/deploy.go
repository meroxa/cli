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
	"os/exec"
	"path"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
)

type Deploy struct {
	flags struct {
		Path string `long:"path" description:"path to the app directory"`
	}

	config config.Config
	logger log.Logger
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
		Example: "meroxa apps deploy # assumes you run it from the app directory" +
			"meroxa apps deploy --path ./my-app",
	}
}

func (d *Deploy) Config(cfg config.Config) {
	d.config = cfg
}

func (d *Deploy) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Deploy) checkRequiredEnvVars() error {
	// TODO: Make sure we could read from either config file or via env vars
	v := os.Getenv(dockerHubUserNameEnv)
	k := os.Getenv(dockerHubAccessTokenEnv)

	if v == "" || k == "" {
		return errors.New("both `DOCKER_HUB_USERNAME` and `DOCKER_HUB_ACCESS_TOKEN` are required to be set to deploy your application")
	}
	return nil
}

func (d *Deploy) getPath() string {
	if d.flags.Path != "" {
		return d.flags.Path
	}
	return "."
}

func (d *Deploy) deployGoApp(ctx context.Context) error {
	appPath := d.getPath()

	appName := path.Base(appPath)
	fqImageName := prependAccount(appName)
	err := buildImage(ctx, d.logger, appPath, fqImageName)
	if err != nil {
		d.logger.Errorf(ctx, "unable to build image; %q\n%s", fqImageName, err)
	}

	err = pushImage(d.logger, fqImageName)
	if err != nil {
		d.logger.Errorf(ctx, "unable to push image; %q\n%s", fqImageName, err)
	}

	err = buildGoApp(ctx, d.logger, appPath, appName, true)
	if err != nil {
		return err
	}

	// deploy data app
	err = deployApp(ctx, d.logger, appPath, appName, fqImageName)
	if err != nil {
		d.logger.Errorf(ctx, "unable to deploy app; %s", err)
	}

	return nil
}

func (d *Deploy) deployJSApp(ctx context.Context) error {
	cmd := exec.Command("npx", "turbine", "deploy", d.getPath()) // nolint:gosec

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		d.logger.Error(ctx, string(stdout))
		return err
	}
	d.logger.Info(ctx, string(stdout))
	return nil
}

func (d *Deploy) Execute(ctx context.Context) error {
	err := d.checkRequiredEnvVars()
	if err != nil {
		return err
	}

	appPath := d.getPath()
	appConfig, err := readConfigFile(appPath)
	if err != nil {
		return err
	}

	lang := appConfig.Language

	if appConfig.Language == "" {
		return fmt.Errorf("`language` should be specified in your app.json")
	}

	switch lang {
	case "go", "golang":
		return d.deployGoApp(ctx)
	case "js", "javascript", nodeJS:
		return d.deployJSApp(ctx)
	default:
		return fmt.Errorf("language %q not supported. Currently, we support \"javascript\" and \"go\"", lang)
	}
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}
