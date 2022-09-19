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
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/volatiletech/null/v8"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine/javascript"
	turbinePY "github.com/meroxa/cli/cmd/meroxa/turbine/python"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

const (
	dockerHubUserNameEnv    = "DOCKER_HUB_USERNAME"
	dockerHubAccessTokenEnv = "DOCKER_HUB_ACCESS_TOKEN" //nolint:gosec

	pollDuration                          = 2 * time.Second
	featureFlagIntermediateRepresentation = "intermediate-representation"
)

type deployApplicationClient interface {
	CreateApplicationV2(ctx context.Context, input *meroxa.CreateApplicationInput) (*meroxa.Application, error)
	CreateDeployment(ctx context.Context, input *meroxa.CreateDeploymentInput) (*meroxa.Deployment, error)
	GetApplication(ctx context.Context, nameOrUUID string) (*meroxa.Application, error)
	GetLatestDeployment(ctx context.Context, appName string) (*meroxa.Deployment, error)
	ListApplications(ctx context.Context) ([]*meroxa.Application, error)
	DeleteApplicationEntities(ctx context.Context, name string) (*http.Response, error)
	CreateBuild(ctx context.Context, input *meroxa.CreateBuildInput) (*meroxa.Build, error)
	CreateSource(ctx context.Context) (*meroxa.Source, error)
	GetBuild(ctx context.Context, uuid string) (*meroxa.Build, error)
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
}

type Deploy struct {
	flags struct {
		Path                     string `long:"path" usage:"Path to the app directory (default is local directory)"`
		DockerHubUserName        string `long:"docker-hub-username" usage:"DockerHub username to use to build and deploy the app" hidden:"true"`         //nolint:lll
		DockerHubAccessToken     string `long:"docker-hub-access-token" usage:"DockerHub access token to use to build and deploy the app" hidden:"true"` //nolint:lll
		Spec                     string `long:"spec" usage:"Deployment specification version to use to build and deploy the app" hidden:"true"`
		SkipCollectionValidation bool   `long:"skip-collection-validation" usage:"Skips unique destination collection and looping validations"` //nolint:lll
	}

	client        deployApplicationClient
	config        config.Config
	logger        log.Logger
	localDeploy   turbine.LocalDeploy
	turbineCLI    turbine.CLI
	appConfig     *turbine.AppConfig
	configAppName string
	appName       string
	gitBranch     string
	path          string
	lang          string
	fnName        string
	tempPath      string // find something more elegant to this
	specVersion   string
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

func (d *Deploy) getPlatformImage(ctx context.Context) (string, error) {
	d.logger.StartSpinner("\t", " Fetching Meroxa Platform source...")

	s, err := d.client.CreateSource(ctx)
	if err != nil {
		d.logger.Errorf(ctx, "\t 𐄂 Unable to fetch source")
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", err
	}
	d.logger.StopSpinnerWithStatus("Platform source fetched", log.Successful)

	err = d.turbineCLI.UploadSource(ctx, d.appName, d.tempPath, s.PutUrl)
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

func (d *Deploy) deployApp(ctx context.Context, imageName, gitSha, specVersion string) error {
	d.logger.Infof(ctx, "Deploying application %q...", d.appName)
	deploymentSpec, err := d.turbineCLI.Deploy(ctx, imageName, d.appName, gitSha, specVersion)
	if err != nil {
		return err
	}
	if specVersion != "" {
		input := &meroxa.CreateDeploymentInput{
			Application: meroxa.EntityIdentifier{Name: null.StringFrom(d.appName)},
			GitSha:      gitSha,
			SpecVersion: null.StringFrom(specVersion),
			Spec:        null.StringFrom(deploymentSpec),
		}
		_, err = d.client.CreateDeployment(ctx, input)
	}
	return err
}

func (d *Deploy) checkPreIRApplication(ctx context.Context) (*meroxa.Application, error) {
	return d.client.GetApplication(ctx, d.appName)
}

// buildApp will call any specifics to the turbine library to prepare a directory that's ready
// to compress, and build, to then later on call the specific command to deploy depending on the language.
func (d *Deploy) buildApp(ctx context.Context) (err error) {
	d.tempPath, err = d.turbineCLI.Build(ctx, d.appName, true)
	return err
}

// getAppImage will check what type of build needs to perform and ultimately will return the image name to use
// when deploying.
func (d *Deploy) getAppImage(ctx context.Context) (string, error) {
	d.logger.StartSpinner("\t", "Checking if application has processes...")
	var fqImageName string

	needsToBuild, err := d.turbineCLI.NeedsToBuild(ctx, d.appName)
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", err
	}

	// If no need to build, return empty imageName which won't be utilized by the deployment process anyway.
	if !needsToBuild {
		d.logger.StopSpinnerWithStatus("No need to create process image...\n", log.Successful)
		return "", nil
	}

	d.logger.StopSpinnerWithStatus("Application processes found. Creating application image...", log.Successful)

	d.localDeploy.TempPath = d.tempPath
	d.localDeploy.Lang = d.lang
	d.localDeploy.AppName = d.appName
	if d.localDeploy.Enabled {
		fqImageName, err = d.localDeploy.GetDockerImageName(ctx, d.logger, d.path, d.appName, d.lang)
		if err != nil {
			return "", err
		}
	} else {
		fqImageName, err = d.getPlatformImage(ctx)
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
	case "go", turbine.GoLang:
		d.lang = turbine.GoLang
	case "js", turbine.JavaScript, turbine.NodeJs:
		d.lang = turbine.JavaScript
	case "py", turbine.Python3, turbine.Python:
		d.lang = turbine.Python
	default:
		return fmt.Errorf("language %q not supported. %s", d.lang, LanguageNotSupportedError)
	}
	return nil
}

func (d *Deploy) validateAppJSON(ctx context.Context) error {
	var err error
	var config turbine.AppConfig

	d.logger.Info(ctx, "Validating your app.json...")
	// validateLocalDeploymentConfig will look for DockerHub credentials to determine whether it's a local deployment or not.
	if err = d.validateLocalDeploymentConfig(); err != nil {
		return err
	}

	if d.path, err = turbine.GetPath(d.flags.Path); err != nil {
		return err
	}

	if d.appConfig == nil {
		d.lang, err = turbine.GetLangFromAppJSON(ctx, d.logger, d.path)
		if err != nil {
			return err
		}
		d.configAppName, err = turbine.GetAppNameFromAppJSON(ctx, d.logger, d.path)
		if err != nil {
			return err
		}
		config, err = turbine.ReadConfigFile(d.path)
		d.appConfig = &config
		if err != nil {
			return err
		}

		if d.gitBranch, err = turbine.GetGitBranch(d.path); err != nil {
			return err
		}
		d.appName = d.prepareAppName(ctx)
	}

	if err = d.validateLanguage(); err != nil {
		return err
	}

	return nil
}

func (d *Deploy) getResourceCheckErrorMessage(ctx context.Context, resources []turbine.ApplicationResource) string {
	var errStr string
	for _, r := range resources {
		resource, err := d.client.GetResourceByNameOrID(ctx, r.Name)
		if err != nil {
			if errStr != "" {
				errStr += "; "
			}

			if strings.Contains(err.Error(), "could not find") {
				errStr += fmt.Sprintf("could not find resource %q", r.Name)
			} else {
				errStr += err.Error()
			}
		} else if resource.Status.State != meroxa.ResourceStateReady {
			if errStr != "" {
				errStr += "; "
			}
			errStr += fmt.Sprintf("resource %q is not ready and usable", r.Name)
		}
	}
	return errStr
}

func (d *Deploy) checkResourceAvailability(ctx context.Context) error {
	resourceCheckMessage := fmt.Sprintf(" Checking resource availability for application %q (%s) before deployment...", d.appName, d.lang)

	d.logger.StartSpinner("\t", resourceCheckMessage)

	resources, err := d.turbineCLI.GetResources(ctx, d.appName)
	if err != nil {
		return fmt.Errorf("unable to read resource definition from app: %s", err.Error())
	}

	if len(resources) == 0 {
		return errors.New("no resources defined in your Turbine app")
	}

	if errStr := d.getResourceCheckErrorMessage(ctx, resources); errStr != "" {
		errStr += ";\n\n ⚠️  Run 'meroxa resources list' to verify that the resource names " +
			"defined in your Turbine app are identical to the resources you have created on the Meroxa Platform before deploying again"
		d.logger.StopSpinnerWithStatus("Resource availability check failed", log.Failed)
		return fmt.Errorf("%s", errStr)
	}

	if !d.flags.SkipCollectionValidation {
		if err := d.validateCollections(ctx, resources); err != nil {
			d.logger.StopSpinnerWithStatus("Resource availability check failed", log.Failed)
			return err
		}
	}

	d.logger.StopSpinnerWithStatus("Can access your Turbine resources", log.Successful)
	return nil
}

func (d *Deploy) prepareDeployment(ctx context.Context) error {
	d.logger.Infof(ctx, "Preparing to deploy application %q...", d.appName)

	// After this point, CLI will package it up and will build it
	err := d.buildApp(ctx)
	if err != nil {
		return err
	}

	// check if resources exist and are ready
	err = d.checkResourceAvailability(ctx)
	if err != nil {
		return err
	}

	d.fnName, err = d.getAppImage(ctx)
	return err
}

// validateSpecVersionDeployment checks, when --spec is specified, whether the version has a valid format
// accepted formats are: semver or "latest".
func (d *Deploy) validateSpecVersionDeployment() error {
	switch d.flags.Spec {
	case "":
		return nil
	case "latest":
		d.specVersion = d.flags.Spec
		return nil
	}

	// check if the version has a valid format
	version, err := semver.NewVersion(d.flags.Spec)
	if err != nil {
		return fmt.Errorf("invalid spec version: %v. You must specify a valid format or use \"latest\"", err)
	}

	d.specVersion = version.String()
	return nil
}

// validateConfig will validate uncommitted changes and git branch.
func (d *Deploy) validateGitConfig(ctx context.Context) error {
	return d.turbineCLI.GitChecks(ctx)
}

type resourceCollectionPair struct {
	collectionName string
	resourceName   string
}

func newResourceCollectionPair(r turbine.ApplicationResource) resourceCollectionPair {
	return resourceCollectionPair{
		collectionName: r.Collection,
		resourceName:   r.Name,
	}
}

func (d *Deploy) validateCollections(ctx context.Context, resources []turbine.ApplicationResource) error {
	var (
		sources      []turbine.ApplicationResource
		destinations = map[resourceCollectionPair]bool{}

		errMessage string
	)
	for _, r := range resources {
		if r.Source && r.Destination {
			errMessage = fmt.Sprintf(
				"%s\n\tApplication resource cannot be used as both a source and destination.",
				errMessage,
			)
		} else if r.Source {
			sources = append(sources, r)
		} else if r.Destination {
			pair := newResourceCollectionPair(r)
			if destinations[pair] {
				errMessage = fmt.Sprintf(
					"%s\n\tApplication resource %q with collection %q cannot be used as a destination more than once.",
					errMessage,
					r.Name,
					r.Collection,
				)
			} else {
				destinations[pair] = true
			}
		}
	}

	apps, err := d.client.ListApplications(ctx)
	if err != nil {
		return err
	}

	errMessage += d.validateNoCollectionLoops(sources, destinations) +
		d.validateDestinationCollectionUnique(apps, destinations)

	if errMessage != "" {
		return fmt.Errorf(
			"⚠️%s\n%s %s",
			errMessage,
			"Please modify your Turbine data application code. Then run `meroxa app deploy` again.",
			"To skip collection validation, run `meroxa app deploy --skip-collection-validation`.",
		)
	}

	return nil
}

// validateNoCollectionLoops ensures source (resource, collection) doesn't equal any destination (resource, collection).
func (d *Deploy) validateNoCollectionLoops(sources []turbine.ApplicationResource, destinations map[resourceCollectionPair]bool) string {
	var errMessage string
	for _, source := range sources {
		if ok := destinations[newResourceCollectionPair(source)]; ok {
			errMessage = fmt.Sprintf(
				"%s\n\tApplication resource %q with collection %q cannot be used as a destination. It is also the source.",
				errMessage,
				source.Name,
				source.Collection)
		}
	}

	return errMessage
}

// validateDestinationCollectionUnique ensures destination (resource, collection) is unique for account.
func (d *Deploy) validateDestinationCollectionUnique(apps []*meroxa.Application, destinations map[resourceCollectionPair]bool) string {
	var errMessage string
	for _, app := range apps {
		for _, r := range app.Resources {
			if r.Collection.Destination.String == "true" &&
				destinations[resourceCollectionPair{
					collectionName: r.Collection.Name.String,
					resourceName:   r.Name.String,
				}] {
				errMessage = fmt.Sprintf(
					"%s\n\tApplication resource %q with collection %q cannot be used as a destination. "+
						"It is also being used as a destination by another application %q.",
					errMessage,
					r.Name.String,
					r.Collection.Name.String,
					app.Name,
				)
			}
		}
	}

	return errMessage
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

func (d *Deploy) prepareAppName(ctx context.Context) string {
	if turbine.GitMainBranch(d.gitBranch) {
		return d.configAppName
	}

	// git branch names can contain a lot of characters that make Docker unhappy
	reg := regexp.MustCompile("[^a-z0-9-_]+")
	sanitizedBranch := reg.ReplaceAllString(strings.ToLower(d.gitBranch), "-")
	appName := fmt.Sprintf("%s-%s", d.configAppName, sanitizedBranch)
	d.logger.Infof(
		ctx,
		"\t%s Feature branch (%s) detected, setting app name to %s...",
		d.logger.SuccessfulCheck(),
		d.gitBranch,
		appName,
	)

	return appName
}

func (d *Deploy) waitForDeployment(ctx context.Context) error {
	startTime := time.Now().UTC()
	logs := []string{}
	for time.Now().UTC().Sub(startTime) < 5*time.Minute {
		var deployment *meroxa.Deployment
		deployment, err := d.client.GetLatestDeployment(ctx, d.appName)
		if err != nil {
			return err
		}
		newLogs := strings.Split(deployment.OutputLog.String, "\n")
		if len(newLogs) > len(logs) {
			for i := len(logs); i < len(newLogs); i++ {
				if len(logs) == 0 && i != 0 {
					d.logger.StopSpinner(newLogs[i-1])
				} else if len(logs) != 0 {
					d.logger.StopSpinner(newLogs[i-1])
				}
				d.logger.StartSpinner("\t", newLogs[i])
			}
			logs = newLogs
		}
		if deployment.Status.State == meroxa.DeploymentStateDeployed {
			l := len(logs)
			d.logger.StopSpinner(logs[l-1])
			return nil
		} else if deployment.Status.State == meroxa.DeploymentStateDeployingError ||
			deployment.Status.State == meroxa.DeploymentStateRollingBackError {
			l := len(logs)
			d.logger.StopSpinnerWithStatus(logs[l-1], log.Failed)
			return fmt.Errorf("failed to deploy Application %q", d.appName)
		}
	}
	return fmt.Errorf("timed out; check `apps logs`")
}

func hasFeatureFlag(f string) bool {
	userFeatureFlags := global.Config.GetStringSlice(global.UserFeatureFlagsEnv)

	for _, v := range userFeatureFlags {
		if v == f {
			return true
		}
	}

	return false
}

//nolint:gocyclo
func (d *Deploy) Execute(ctx context.Context) error {
	var app *meroxa.Application

	if err := d.validateAppJSON(ctx); err != nil {
		return err
	}

	switch d.lang {
	case "go", turbine.GoLang:
		if d.turbineCLI == nil {
			d.turbineCLI = turbineGo.New(d.logger, d.path)
		}
	case "js", turbine.JavaScript, turbine.NodeJs:
		if d.turbineCLI == nil {
			d.turbineCLI = turbineJS.New(d.logger, d.path)
		}
	case "py", turbine.Python3, turbine.Python:
		if d.turbineCLI == nil {
			d.turbineCLI = turbinePY.New(d.logger, d.path)
		}
	default:
		return fmt.Errorf("language %q not supported. %s", d.lang, LanguageNotSupportedError)
	}

	if err := d.validateGitConfig(ctx); err != nil {
		return err
	}

	gitSha, err := d.turbineCLI.GetGitSha(ctx)
	if err != nil {
		return err
	}

	if d.specVersion == "" {
		// ⚠️ This is only until we re-deploy applications applying only the changes made
		if err = d.tearDownExistingResources(ctx); err != nil {
			return err
		}
	} else {
		// Intermediate Representation Workflow
		if err = d.validateSpecVersionDeployment(); err != nil {
			return err
		}
		if !hasFeatureFlag(featureFlagIntermediateRepresentation) {
			return fmt.Errorf("user is not authorized for deploying with --spec")
		}

		app, err = d.client.CreateApplicationV2(ctx, &meroxa.CreateApplicationInput{
			Name:     d.appName,
			Language: d.lang,
			GitSha:   null.StringFrom(gitSha)})
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}

	err = d.prepareDeployment(ctx)
	defer d.turbineCLI.CleanUpBinaries(ctx, d.appName)
	if err != nil {
		return err
	}

	if err = d.deployApp(ctx, d.fnName, gitSha, d.specVersion); err != nil {
		return err
	}

	if d.specVersion == "" {
		app, err = d.checkPreIRApplication(ctx)
	} else {
		err = d.waitForDeployment(ctx)
	}
	if err != nil {
		return err
	}

	dashboardURL := fmt.Sprintf("https://dashboard.meroxa.io/apps/%s/detail", d.appName)
	output := fmt.Sprintf("\t%s Application %q successfully created!\n\n  ✨ To visualize your application visit %s",
		d.logger.SuccessfulCheck(), d.appName, dashboardURL)
	d.logger.Info(ctx, output)
	d.logger.JSON(ctx, app)
	return nil
}
