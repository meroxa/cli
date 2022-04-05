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
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine_cli/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine_cli/javascript"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	turbine "github.com/meroxa/turbine-go/deploy"
)

const (
	dockerHubUserNameEnv    = "DOCKER_HUB_USERNAME"
	dockerHubAccessTokenEnv = "DOCKER_HUB_ACCESS_TOKEN" // nolint:gosec
	pollDuration            = 2 * time.Second
)

type deployApplicationClient interface {
	CreateApplication(ctx context.Context, input *meroxa.CreateApplicationInput) (*meroxa.Application, error)
	GetApplication(ctx context.Context, nameOrUUID string) (*meroxa.Application, error)
	CreateBuild(ctx context.Context, input *meroxa.CreateBuildInput) (*meroxa.Build, error)
	CreateSource(ctx context.Context) (*meroxa.Source, error)
	GetBuild(ctx context.Context, uuid string) (*meroxa.Build, error)
}

type Deploy struct {
	flags struct {
		Path                 string `long:"path" description:"path to the app directory (default is local directory)"`
		DockerHubUserName    string `long:"docker-hub-username" description:"DockerHub username to use to build and deploy the app" hidden:"true"`         //nolint:lll
		DockerHubAccessToken string `long:"docker-hub-access-token" description:"DockerHub access token to use to build and deploy the app" hidden:"true"` //nolint:lll
	}

	client   deployApplicationClient
	config   config.Config
	logger   log.Logger
	appName  string
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

// TODO: Move this to each turbine library.
func (d *Deploy) createApplication(ctx context.Context, pipelineUUID, gitSha string) error {
	appName, err := turbineCLI.GetAppNameFromAppJSON(d.path)
	if err != nil {
		return err
	}

	input := meroxa.CreateApplicationInput{
		Name:     appName,
		Language: d.lang,
		GitSha:   gitSha,
		Pipeline: meroxa.EntityIdentifier{UUID: null.StringFrom(pipelineUUID)},
	}
	d.logger.Infof(ctx, "Creating application %q with language %q...", input.Name, d.lang)

	res, err := d.client.CreateApplication(ctx, &input)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			// Double check that the created application has the expected pipeline.
			var app *meroxa.Application
			app, err = d.client.GetApplication(ctx, appName)
			if err != nil {
				return err
			}
			if app.Pipeline.UUID.String != pipelineUUID {
				return fmt.Errorf("unable to finish creating the %s Application because its entities are in an"+
					" unrecoverable state; try deleting and re-deploying", appName)
			}
		}
		return err
	}

	d.logger.Infof(ctx, "Application %q successfully created!", res.Name)
	d.logger.JSON(ctx, res)
	return nil
}

// uploadSource creates first a Dockerfile to then, package the entire folder which will be later uploaded
// this should ignore .git files and fixtures/.
func (d *Deploy) uploadSource(ctx context.Context, appPath, url string) error {
	// Before creating a .tar.zip, we make sure it contains a Dockerfile.
	err := turbine.CreateDockerfile(appPath)
	if err != nil {
		return err
	}

	dFile := fmt.Sprintf("turbine-%s.tar.gz", d.appName)

	var buf bytes.Buffer
	d.logger.Infof(ctx, "Packaging application located at %q...", appPath)
	err = turbineCLI.CreateTarAndZipFile(appPath, &buf)
	if err != nil {
		return err
	}

	fileToWrite, err := os.OpenFile(dFile, os.O_CREATE|os.O_RDWR, os.FileMode(0777)) //nolint:gomnd
	defer func(fileToWrite *os.File) {
		err = fileToWrite.Close()
		if err != nil {
			panic(err.Error())
		}
	}(fileToWrite)

	if err != nil {
		return err
	}
	if _, err = io.Copy(fileToWrite, &buf); err != nil {
		return err
	}

	// We clean up Dockerfile as last step
	err = os.Remove(filepath.Join(appPath, "Dockerfile"))
	if err != nil {
		return err
	}

	err = d.uploadFile(ctx, dFile, url)
	if err != nil {
		return err
	}

	// remove .tar.gz file
	return os.Remove(dFile)
}

func (d *Deploy) uploadFile(ctx context.Context, filePath, url string) error {
	d.logger.Info(ctx, "Uploading file to our build service...")
	fh, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(fh *os.File) {
		err = fh.Close()
		if err != nil {
			d.logger.Warn(ctx, err.Error())
		}
	}(fh)

	req, err := http.NewRequestWithContext(ctx, "PUT", url, fh)
	if err != nil {
		return err
	}

	fi, err := fh.Stat()
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	req.ContentLength = fi.Size()

	client := &http.Client{}
	res, err := client.Do(req) //nolint:bodyclose
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		d.logger.Infof(ctx, "Uploaded!")
		if err != nil {
			d.logger.Error(ctx, err.Error())
		}
	}(res.Body)

	return nil
}

func (d *Deploy) getPlatformImage(ctx context.Context, appPath string) (string, error) {
	s, err := d.client.CreateSource(ctx)
	if err != nil {
		return "", err
	}

	err = d.uploadSource(ctx, appPath, s.PutUrl)
	if err != nil {
		return "", err
	}

	sourceBlob := meroxa.SourceBlob{Url: s.GetUrl}
	buildInput := &meroxa.CreateBuildInput{SourceBlob: sourceBlob}

	build, err := d.client.CreateBuild(ctx, buildInput)
	if err != nil {
		return "", err
	}

	fmt.Printf("Getting status for build: %s ", build.Uuid)
	for {
		fmt.Printf(".")
		b, err := d.client.GetBuild(ctx, build.Uuid)
		if err != nil {
			return "", err
		}

		switch b.Status.State {
		case "error":
			return "", fmt.Errorf("build with uuid %q errored ", b.Uuid)
		case "complete":
			fmt.Println("\nImage built! ")
			return build.Image, nil
		}
		time.Sleep(pollDuration)
	}
}

// Deploy takes care of all the necessary steps to deploy a Turbine application
//	1. Build binary // different for jS
//	2. Build image // common
//	3. Push image // common
//	4. Run Turbine deploy // different
func (d *Deploy) deploy(ctx context.Context, appPath string, l log.Logger) (string, error) {
	var fqImageName string
	d.appName = path.Base(appPath)

	err := turbineGo.BuildBinary(ctx, l, appPath, d.appName, true)
	if err != nil {
		return "", err
	}

	var ok bool
	// check for image instances
	if ok, err = turbineGo.NeedsToBuild(appPath, d.appName); ok {
		if err != nil {
			l.Errorf(ctx, err.Error())
			return "", err
		}

		if d.goDeploy.LocalDeployment {
			fqImageName, err = d.goDeploy.GetDockerImageName(ctx, l, appPath, d.appName)
			if err != nil {
				return "", err
			}
		} else {
			fqImageName, err = d.getPlatformImage(ctx, appPath)
			if err != nil {
				return "", err
			}
		}
	}

	// creates all resources
	output, err := turbineGo.RunDeployApp(ctx, l, appPath, d.appName, fqImageName)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (d *Deploy) Execute(ctx context.Context) error {
	// validateLocalDeploymentConfig will look for DockerHub credentials to determine whether it's a local deployment or not.
	err := d.validateLocalDeploymentConfig()
	if err != nil {
		return err
	}
	var deployOutput string

	d.path, err = turbineCLI.GetPath(d.flags.Path)
	if err != nil {
		return err
	}
	d.lang, err = turbineCLI.GetLangFromAppJSON(d.path)
	if err != nil {
		return err
	}

	err = turbineCLI.GitChecks(ctx, d.logger, d.path)
	if err != nil {
		return err
	}

	err = turbineCLI.ValidateBranch(d.path)
	if err != nil {
		return err
	}

	// 1. set up the app structure (CLI does this for any language)
	// 2. *depending on the language* call something to create the dockerfile <=
	// 3. CLI would handle:
	//	3.1 creating the tar.zip,
	//	3.2 post /sources
	// 	3.3 uploading the tar.zip
	//  3.4 post /builds
	// 4. CLI would call (depending on language) the deploy script <=
	switch d.lang {
	case GoLang:
		deployOutput, err = d.deploy(ctx, d.path, d.logger)
	case "js", JavaScript, NodeJs:
		err = turbineJS.Deploy(ctx, d.path, d.logger)
	default:
		return fmt.Errorf("language %q not supported. %s", d.lang, LanguageNotSupportedError)
	}
	if err != nil {
		return err
	}

	pipelineUUID := turbineCLI.GetPipelineUUID(deployOutput)
	gitSha, err := turbineCLI.GetGitSha(d.path)
	if err != nil {
		return err
	}

	return d.createApplication(ctx, pipelineUUID, gitSha)
}
