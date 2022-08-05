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
	"path/filepath"
	"strings"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine_cli/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine_cli/javascript"
	turbinePY "github.com/meroxa/cli/cmd/meroxa/turbine_cli/python"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	turbine "github.com/meroxa/turbine-go/deploy"
)

const (
	dockerHubUserNameEnv    = "DOCKER_HUB_USERNAME"
	dockerHubAccessTokenEnv = "DOCKER_HUB_ACCESS_TOKEN" //nolint:gosec
	pollDuration            = 2 * time.Second
)

type deployApplicationClient interface {
	CreateApplication(ctx context.Context, input *meroxa.CreateApplicationInput) (*meroxa.Application, error)
	GetApplication(ctx context.Context, nameOrUUID string) (*meroxa.Application, error)
	ListApplications(ctx context.Context) ([]*meroxa.Application, error)
	DeleteApplicationEntities(ctx context.Context, name string) (*http.Response, error)
	CreateBuild(ctx context.Context, input *meroxa.CreateBuildInput) (*meroxa.Build, error)
	CreateSource(ctx context.Context) (*meroxa.Source, error)
	GetBuild(ctx context.Context, uuid string) (*meroxa.Build, error)
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
}

type Deploy struct {
	flags struct {
		Path                 string `long:"path" description:"path to the app directory (default is local directory)"`
		DockerHubUserName    string `long:"docker-hub-username" description:"DockerHub username to use to build and deploy the app" hidden:"true"`         //nolint:lll
		DockerHubAccessToken string `long:"docker-hub-access-token" description:"DockerHub access token to use to build and deploy the app" hidden:"true"` //nolint:lll
	}

	client      deployApplicationClient
	config      config.Config
	logger      log.Logger
	appName     string
	path        string
	lang        string
	localDeploy turbineCLI.LocalDeploy
	fnName      string
	tempPath    string // find something more elegant to this
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
		Short: "Deploy a Turbine Data Application",
		Long: `This command will deploy the application specified in '--path'
(or current working directory if not specified) to our Meroxa Platform.
If deployment was successful, you should expect an application you'll be able to fully manage
`,
		Example: `meroxa apps deploy # assumes you run it from the app directory
meroxa apps deploy --path ./my-app
`,
		Beta: true,
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
		d.localDeploy.DockerHubUserNameEnv = d.flags.DockerHubUserName
		if d.flags.DockerHubAccessToken == "" {
			return errors.New("--docker-hub-access-token is required when --docker-hub-username is present")
		}
	}

	if d.flags.DockerHubAccessToken != "" {
		d.localDeploy.DockerHubAccessTokenEnv = d.flags.DockerHubAccessToken
		if d.flags.DockerHubUserName == "" {
			return errors.New("--docker-hub-username is required when --docker-hub-access-token is present")
		}
	}
	return nil
}

func (d *Deploy) validateDockerHubEnvVars() error {
	if d.getDockerHubUserNameEnv() != "" {
		d.localDeploy.DockerHubUserNameEnv = d.getDockerHubUserNameEnv()
		if d.getDockerHubAccessTokenEnv() == "" {
			return fmt.Errorf("%s is required when %s is present", dockerHubAccessTokenEnv, dockerHubUserNameEnv)
		}
	}
	if d.getDockerHubAccessTokenEnv() != "" {
		d.localDeploy.DockerHubAccessTokenEnv = d.getDockerHubAccessTokenEnv()
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
	if d.localDeploy.DockerHubUserNameEnv != "" && d.localDeploy.DockerHubAccessTokenEnv != "" {
		d.localDeploy.Enabled = true
	}
	return nil
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}

// uploadSource creates first a Dockerfile to then, package the entire folder which will be later uploaded
// this should ignore .git files and fixtures/.
func (d *Deploy) uploadSource(ctx context.Context, appPath, url string) error {
	var err error

	if d.lang == GoLang {
		d.logger.StartSpinner("\t", fmt.Sprintf("Creating Dockerfile before uploading source in %s", appPath))
		err = turbine.CreateDockerfile("", appPath)
		if err != nil {
			return err
		}
		d.logger.StopSpinnerWithStatus("Dockerfile created", log.Successful)
	}

	dFile := fmt.Sprintf("turbine-%s.tar.gz", d.appName)

	var buf bytes.Buffer

	if d.lang == JavaScript || d.lang == Python {
		appPath = d.tempPath
	}

	err = turbineCLI.CreateTarAndZipFile(appPath, &buf)
	if err != nil {
		return err
	}

	d.logger.StartSpinner("\t", fmt.Sprintf(" Creating %q in %q to upload to our build service...", appPath, dFile))

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
	d.logger.StopSpinnerWithStatus(fmt.Sprintf("%q successfully created in %q", dFile, appPath), log.Successful)

	if d.lang == GoLang {
		d.logger.StartSpinner("\t", fmt.Sprintf("Removing Dockerfile created for your application in %s...", appPath))
		// We clean up Dockerfile as last step
		err = os.Remove(filepath.Join(appPath, "Dockerfile"))
		if err != nil {
			return err
		}
		d.logger.StopSpinnerWithStatus("Dockerfile removed", log.Successful)
	}

	err = d.uploadFile(ctx, dFile, url)
	if err != nil {
		return err
	}

	if d.lang == Python {
		var output string
		output, err = turbinePY.CleanUpApp(appPath)
		if err != nil {
			fmt.Printf("warning: failed to clean up app at %s: %v %s\n", appPath, err, output)
		}
	}
	// remove .tar.gz file
	d.logger.StartSpinner("\t", fmt.Sprintf(" Removing %q...", dFile))
	err = os.Remove(dFile)
	if err != nil {
		d.logger.StopSpinnerWithStatus(fmt.Sprintf("\t Something went wrong trying to remove %q", dFile), log.Failed)
		return err
	}
	d.logger.StopSpinnerWithStatus(fmt.Sprintf("%q removed", dFile), log.Successful)
	return nil
}

func (d *Deploy) uploadFile(ctx context.Context, filePath, url string) error {
	d.logger.StartSpinner("\t", " Uploading source...")

	fh, err := os.Open(filePath)
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}
	defer func(fh *os.File) {
		fh.Close()
	}(fh)

	req, err := http.NewRequestWithContext(ctx, "PUT", url, fh)
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	fi, err := fh.Stat()
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	req.ContentLength = fi.Size()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	d.logger.StopSpinnerWithStatus("Source uploaded", log.Successful)
	return nil
}

func (d *Deploy) getPlatformImage(ctx context.Context, appPath string) (string, error) {
	d.logger.StartSpinner("\t", " Fetching Meroxa Platform source...")

	s, err := d.client.CreateSource(ctx)
	if err != nil {
		d.logger.Errorf(ctx, "\t 𐄂 Unable to fetch source")
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", err
	}
	d.logger.StopSpinnerWithStatus("Platform source fetched", log.Successful)

	err = d.uploadSource(ctx, appPath, s.PutUrl)
	if err != nil {
		return "", err
	}

	sourceBlob := meroxa.SourceBlob{Url: s.GetUrl}
	buildInput := &meroxa.CreateBuildInput{SourceBlob: sourceBlob}

	build, err := d.client.CreateBuild(ctx, buildInput)
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", err
	}
	d.logger.StartSpinner("\t", fmt.Sprintf(" Building Meroxa Process image (%q)...", build.Uuid))

	for {
		b, err := d.client.GetBuild(ctx, build.Uuid)
		if err != nil {
			d.logger.StopSpinnerWithStatus("\t", log.Failed)
			return "", err
		}

		switch b.Status.State {
		case "error":
			msg := fmt.Sprintf("build with uuid %q errored\nRun `meroxa build logs %s` for more information", b.Uuid, b.Uuid)
			d.logger.StopSpinnerWithStatus(msg, log.Failed)
			return "", fmt.Errorf("build with uuid %q errored", b.Uuid)
		case "complete":
			d.logger.StopSpinnerWithStatus(fmt.Sprintf("Successfully built process image (%q)\n", build.Uuid), log.Successful)
			return build.Image, nil
		}
		time.Sleep(pollDuration)
	}
}

func (d *Deploy) deployApp(ctx context.Context, imageName, gitSha string) error {
	var err error

	d.logger.StartSpinner("\t", fmt.Sprintf(" Deploying application %q...", d.appName))
	switch d.lang {
	case GoLang:
		err = turbineGo.RunDeployApp(ctx, d.logger, d.path, imageName, d.appName, gitSha)
	case JavaScript:
		err = turbineJS.RunDeployApp(ctx, d.logger, d.path, imageName, gitSha)
	case Python:
		err = turbinePY.RunDeployApp(ctx, d.logger, d.path, imageName, gitSha)
	}
	if err != nil {
		d.logger.StopSpinnerWithStatus("Deployment failed\n\n", log.Failed)
		return err
	}

	app, err := d.client.GetApplication(ctx, d.appName)
	if err != nil {
		d.logger.StopSpinnerWithStatus("Deployment failed to create Application\n\n", log.Failed)
		return err
	}
	d.logger.StopSpinnerWithStatus("Deploy complete", log.Successful)

	dashboardURL := fmt.Sprintf("https://dashboard.meroxa.io/apps/%s/detail", app.UUID)
	output := fmt.Sprintf("\t%s Application %q successfully created!\n\n  ✨ To visualize your application visit %s",
		d.logger.SuccessfulCheck(), app.Name, dashboardURL)
	d.logger.StopSpinner(output)
	d.logger.JSON(ctx, app)
	return nil
}

// buildApp will call any specifics to the turbine library to prepare a directory that's ready
// to compress, and build, to then later on call the specific command to deploy depending on the language.
func (d *Deploy) buildApp(ctx context.Context) error {
	var err error

	// Without the " " at the beginning of `suffix`, spinner looks next to word (only on this occurrence)
	d.logger.StartSpinner("\t", " Building application...")

	switch d.lang {
	case GoLang:
		err = turbineGo.BuildBinary(ctx, d.logger, d.path, d.appName, true)
	case JavaScript:
		d.tempPath, err = turbineJS.BuildApp(d.path)
	case Python:
		// Dockerfile will already exist
		d.tempPath, err = turbinePY.BuildApp(d.path)
	}
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}
	d.logger.StopSpinnerWithStatus("Application built", log.Successful)
	return nil
}

// getAppImage will check what type of build needs to perform and ultimately will return the image name to use
// when deploying.
func (d *Deploy) getAppImage(ctx context.Context) (string, error) {
	d.logger.StartSpinner("\t", "Checking if application has processes...")
	var fqImageName string
	var needsToBuild bool
	var err error

	switch d.lang {
	case GoLang:
		needsToBuild, err = turbineGo.NeedsToBuild(d.path, d.appName)
	case JavaScript:
		needsToBuild, err = turbineJS.NeedsToBuild(d.path)
	case Python:
		needsToBuild, err = turbinePY.NeedsToBuild(d.path)
	}
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", err
	}

	// If no need to build, return empty imageName which won't be utilized by the deploy process anyways
	if !needsToBuild {
		d.logger.StopSpinnerWithStatus("No need to create process image...\n", log.Successful)
		return "", nil
	}

	d.logger.StopSpinnerWithStatus("Application processes found. Creating application image...", log.Successful)

	d.localDeploy.TempPath = d.tempPath
	d.localDeploy.Lang = d.lang
	if d.localDeploy.Enabled {
		fqImageName, err = d.localDeploy.GetDockerImageName(ctx, d.logger, d.path, d.appName, d.lang)
		if err != nil {
			return "", err
		}
	} else {
		fqImageName, err = d.getPlatformImage(ctx, d.path)
		if err != nil {
			return "", err
		}
	}

	return fqImageName, nil
}

// validateLanguage stops execution of the deployment in case language is not supported.
// It also consolidates lang used in API in case user specified a supported language using an unexpected description.
func (d *Deploy) validateLanguage() error {
	switch d.lang {
	case "go", GoLang:
		d.lang = GoLang
	case "js", JavaScript, NodeJs:
		d.lang = JavaScript
	case "py", Python3, Python:
		d.lang = Python
	default:
		return fmt.Errorf("language %q not supported. %s", d.lang, LanguageNotSupportedError)
	}
	return nil
}

func (d *Deploy) validate(ctx context.Context) error {
	// validateLocalDeploymentConfig will look for DockerHub credentials to determine whether it's a local deployment or not.
	err := d.validateLocalDeploymentConfig()
	if err != nil {
		return err
	}

	d.path, err = turbineCLI.GetPath(d.flags.Path)
	if err != nil {
		return err
	}

	d.lang, err = turbineCLI.GetLangFromAppJSON(ctx, d.logger, d.path)
	if err != nil {
		return err
	}

	err = d.validateLanguage()
	if err != nil {
		return err
	}

	d.appName, err = turbineCLI.GetAppNameFromAppJSON(ctx, d.logger, d.path)
	if err != nil {
		return err
	}

	err = turbineCLI.GitChecks(ctx, d.logger, d.path)
	if err != nil {
		return err
	}

	return turbineCLI.ValidateBranch(ctx, d.logger, d.path)
}

func (d *Deploy) getResourceCheckErrorMessage(ctx context.Context, resourceNames []string) string {
	var errStr string
	for _, name := range resourceNames {
		resource, err := d.client.GetResourceByNameOrID(ctx, name)
		if err != nil {
			if errStr != "" {
				errStr += "; "
			}

			if strings.Contains(err.Error(), "could not find") {
				errStr += fmt.Sprintf("could not find resource %q", name)
			} else {
				errStr += err.Error()
			}
		} else if resource.Status.State != meroxa.ResourceStateReady {
			if errStr != "" {
				errStr += "; "
			}
			errStr += fmt.Sprintf("resource %q is not ready and usable", resource.Name)
		}
	}
	return errStr
}

func (d *Deploy) checkResourceAvailability(ctx context.Context) error {
	resourceCheckMessage := fmt.Sprintf(" Checking resource availability for application %q (%s) before deployment...", d.appName, d.lang)

	d.logger.StartSpinner("\t", resourceCheckMessage)

	var resourceNames []string
	var err error

	switch d.lang {
	case GoLang:
		resourceNames, err = turbineGo.GetResourceNames(ctx, d.logger, d.path, d.appName)
	case JavaScript:
		resourceNames, err = turbineJS.GetResourceNames(ctx, d.logger, d.path, d.appName)
	case Python:
		resourceNames, err = turbinePY.GetResourceNames(ctx, d.logger, d.path, d.appName)
	}

	if err != nil {
		return fmt.Errorf("unable to read resource definition from app: %s", err.Error())
	}

	if len(resourceNames) == 0 {
		return errors.New("no resources defined in your Turbine app")
	}

	errStr := d.getResourceCheckErrorMessage(ctx, resourceNames)

	if errStr != "" {
		errStr += ";\n\n ⚠️  Run 'meroxa resources list' to verify that the resource names " +
			"defined in your Turbine app are identical to the resources you have created on the Meroxa Platform before deploying again"
		d.logger.StopSpinnerWithStatus("Resource availability check failed", log.Failed)
		return fmt.Errorf("%s", errStr)
	}

	d.logger.StopSpinnerWithStatus("Can access to your Turbine resources", log.Successful)
	return nil
}

func (d *Deploy) prepareAppForDeployment(ctx context.Context) error {
	d.logger.Infof(ctx, "Deploying application %q...", d.appName)

	// After this point, CLI will package it up and will build it
	err := d.buildApp(ctx)
	if err != nil {
		return err
	}

	// check that resources exist and are ready
	err = d.checkResourceAvailability(ctx)
	if err != nil {
		return err
	}

	d.fnName, err = d.getAppImage(ctx)
	return err
}

func (d *Deploy) rmBinary() {
	if d.lang == GoLang {
		localBinary := filepath.Join(d.path, d.appName)
		err := os.Remove(localBinary)
		if err != nil {
			fmt.Printf("warning: failed to clean up %s\n", localBinary)
		}

		crossCompiledBinary := filepath.Join(d.path, d.appName) + ".cross"
		err = os.Remove(crossCompiledBinary)
		if err != nil {
			fmt.Printf("warning: failed to clean up %s\n", crossCompiledBinary)
		}
	}
}

// tearDownExistingResources will only delete the application and its associated entities if it hasn't been created
// or whether it's in a non-running state.
func (d *Deploy) tearDownExistingResources(ctx context.Context) error {
	app, _ := d.client.GetApplication(ctx, d.appName)

	if app != nil && app.Status.State == meroxa.ApplicationStateRunning {
		appIsReady := fmt.Sprintf("application %q is already %s", d.appName, app.Status.State)
		msg := fmt.Sprintf("%s\n\t. Use `meroxa apps remove %s` if you want to redeploy to this application", appIsReady, d.appName)
		return errors.New(msg)
	}
	resp, _ := d.client.DeleteApplicationEntities(ctx, d.appName)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	return nil
}

func (d *Deploy) Execute(ctx context.Context) error {
	d.logger.Info(ctx, "Validating your app.json...")
	err := d.validate(ctx)
	if err != nil {
		return err
	}

	// ⚠️ This is only until we re-deploy applications applying only the changes made
	err = d.tearDownExistingResources(ctx)
	if err != nil {
		return err
	}

	err = d.prepareAppForDeployment(ctx)
	defer d.rmBinary()
	if err != nil {
		return err
	}

	gitSha, err := turbineCLI.GetGitSha(d.path)
	if err != nil {
		return err
	}

	return d.deployApp(ctx, d.fnName, gitSha)
}
