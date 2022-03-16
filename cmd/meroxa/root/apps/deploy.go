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
	"regexp"
	"strings"

	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine_cli/golang"

	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine_cli/javascript"

	"github.com/volatiletech/null/v8"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

const (
	dockerHubUserNameEnv    = "DOCKER_HUB_USERNAME"
	dockerHubAccessTokenEnv = "DOCKER_HUB_ACCESS_TOKEN" // nolint:gosec
)

type createApplicationClient interface {
	CreateApplication(ctx context.Context, input *meroxa.CreateApplicationInput) (*meroxa.Application, error)
}

type Deploy struct {
	flags struct {
		Path                 string `long:"path" description:"path to the app directory (default is local directory)"`
		DockerHubUserName    string `long:"docker-hub-username" description:"DockerHub username to use to build and deploy the application image"`
		DockerHubAccessToken string `long:"docker-hub-access-token" description:"DockerHub access token to use to build and deploy the application image"`
	}

	client   createApplicationClient
	config   config.Config
	logger   log.Logger
	path     string
	lang     string
	goDeploy turbineGo.Deploy
}

var (
	_ builder.CommandWithClient  = (*Deploy)(nil)
	_ builder.CommandWithConfig  = (*Deploy)(nil)
	_ builder.CommandWithDocs    = (*Deploy)(nil)
	_ builder.CommandWithExecute = (*Deploy)(nil)
	_ builder.CommandWithFlags   = (*Deploy)(nil)
	_ builder.CommandWithLogger  = (*Deploy)(nil)
)

func (*Deploy) Usage() string {
	return "deploy"
}

func (*Deploy) Docs() builder.Docs {
	return builder.Docs{
		Short: "Deploy a Meroxa Data Application",
		Long: "This command will deploy the application specified in `--path`" +
			"(or current working directory if not specified) to our Meroxa Platform." +
			"If deployment was successful, you should expect an application you'll be able to fully manage",
		Example: "meroxa apps deploy # assumes you run it from the app directory\n" +
			"meroxa apps deploy --path ./my-app",
	}
}

func (d *Deploy) Config(cfg config.Config) {
	d.config = cfg
}

func (d *Deploy) Client(client meroxa.Client) {
	d.client = client
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

func (d *Deploy) validateDockerHubFlags() error {
	if d.flags.DockerHubUserName != "" {
		d.goDeploy.DockerHubUserNameEnv = d.flags.DockerHubUserName
		if d.flags.DockerHubAccessToken == "" {
			return errors.New("--docker-hub-access-token is required when --docker-hub-username is present")
		}
	}

	if d.flags.DockerHubAccessToken != "" {
		d.goDeploy.DockerHubAccessTokenEnv = d.flags.DockerHubAccessToken
		if d.flags.DockerHubUserName == "" {
			return errors.New("--docker-hub-username is required when --docker-hub-access-token is present")
		}
	}
	return nil
}

func (d *Deploy) validateDockerHubEnvVars() error {
	if d.getDockerHubUserNameEnv() != "" {
		d.goDeploy.DockerHubUserNameEnv = d.getDockerHubUserNameEnv()
		if d.getDockerHubAccessTokenEnv() == "" {
			return fmt.Errorf("%s is required when %s is present", dockerHubAccessTokenEnv, dockerHubUserNameEnv)
		}
	}
	if d.getDockerHubAccessTokenEnv() != "" {
		d.goDeploy.DockerHubAccessTokenEnv = d.getDockerHubAccessTokenEnv()
		if d.getDockerHubUserNameEnv() == "" {
			return fmt.Errorf("%s is required when %s is present", dockerHubUserNameEnv, dockerHubAccessTokenEnv)
		}
	}
	return nil
}

func (d *Deploy) validateLocalDeploymentConfig() error {
	// Check if users had an old configuration where DockerHub credentials were set as environment variables
	err := d.validateDockerHubEnvVars()
	if err != nil {
		return err
	}

	// Check if users are making use of DockerHub flags
	err = d.validateDockerHubFlags()
	if err != nil {
		return err
	}

	// If there are DockerHub credentials, we consider it a local deployment
	if d.goDeploy.DockerHubUserNameEnv != "" && d.goDeploy.DockerHubAccessTokenEnv != "" {
		d.goDeploy.LocalDeployment = true
	}
	return nil
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}

// TODO: Move this to each turbine library
func (d *Deploy) createApplication(ctx context.Context, pipelineUUID string) error {
	appName, err := turbineCLI.GetAppNameFromAppJSON(d.path)
	if err != nil {
		return err
	}

	input := meroxa.CreateApplicationInput{
		Name:     appName,
		Language: d.lang,
		GitSha:   "hardcoded",
		Pipeline: meroxa.EntityIdentifier{UUID: null.StringFrom(pipelineUUID)},
	}
	d.logger.Infof(ctx, "Creating application %q with language %q...", input.Name, d.lang)

	res, err := d.client.CreateApplication(ctx, &input)
	if err != nil {
		return err
	}

	d.logger.Infof(ctx, "Application %q successfully created!", res.Name)
	d.logger.JSON(ctx, res)
	return nil
}

// gitChecks prints warnings about uncommitted tracked and untracked files.
func (d *Deploy) gitChecks(ctx context.Context) error {
	// temporarily switching to the app's directory
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(d.path)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "status", "--porcelain=v2")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	all := string(output)
	lines := strings.Split(strings.TrimSpace(all), "\n")
	if len(lines) > 0 && lines[0] != "" {
		cmd = exec.Command("git", "status")
		output, err = cmd.Output()
		if err != nil {
			return err
		}
		d.logger.Error(ctx, string(output))
		err = os.Chdir(pwd)
		if err != nil {
			return err
		}
		return fmt.Errorf("unable to proceed with deployment because of uncommitted changes")
	}

	return os.Chdir(pwd)
}

// getPipelineUUID parses the deploy output when it was successful to determine the pipeline UUID to create
func getPipelineUUID(output string) string {
	// Example output:
	// 2022/03/16 13:21:36 pipeline created: "turbine-pipeline-simple" ("049760a8-a3d2-44d9-b326-0614c09a3f3e")
	re := regexp.MustCompile(`pipeline created:."[a-zA-Z]+-[a-zA-Z]+-[a-zA-Z]+".(\([^)]*\))`)
	res := re.FindStringSubmatch(output)[1]
	res = strings.Trim(res, "()\"")
	return res
}

func (d *Deploy) Execute(ctx context.Context) error {
	var deployOuput string
	// validateLocalDeploymentConfig will look for DockerHub credentials to determine whether it's a local deployment or not
	err := d.validateLocalDeploymentConfig()
	if err != nil {
		return err
	}

	d.path = turbineCLI.GetPath(d.flags.Path)
	d.lang, err = turbineCLI.GetLangFromAppJSON(d.path)
	if err != nil {
		return err
	}

	err = d.gitChecks(ctx)
	if err != nil {
		return err
	}

	switch d.lang {
	case GoLang:
		// The only reason Deploy is scoped this other way is, so we can have the Docker Credentials
		// Maybe that function should take care of checking type of deployment, only passing flags
		// and environment variables
		// err = turbineGo.Deploy(ctx, d.path, d.logger)
		deployOuput, err = d.goDeploy.Deploy(ctx, d.path, d.logger)
	case "js", JavaScript, NodeJs:
		err = turbineJS.Deploy(ctx, d.path, d.logger)
	default:
		return fmt.Errorf("language %q not supported. %s", d.lang, LanguageNotSupportedError)
	}
	if err != nil {
		return err
	}

	pipelineUUID := getPipelineUUID(deployOuput)

	return d.createApplication(ctx, pipelineUUID)
}
