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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/turbine-core/pkg/ir"
)

const (
	platformBuildPollDuration  = 2 * time.Second
	minutesToWaitForDeployment = 5
	intervalCheckForDeployment = 500 * time.Millisecond
)

type deployApplicationClient interface {
	CreateApplicationV2(ctx context.Context, input *meroxa.CreateApplicationInput) (*meroxa.Application, error)
	CreateDeployment(ctx context.Context, input *meroxa.CreateDeploymentInput) (*meroxa.Deployment, error)
	GetApplication(ctx context.Context, nameOrUUID string) (*meroxa.Application, error)
	GetEnvironment(ctx context.Context, nameOrUUID string) (*meroxa.Environment, error)
	DeleteApplicationEntities(ctx context.Context, nameOrUUID string) (*http.Response, error)
	GetDeployment(ctx context.Context, appName string, depUUID string) (*meroxa.Deployment, error)
	ListApplications(ctx context.Context) ([]*meroxa.Application, error)
	CreateBuild(ctx context.Context, input *meroxa.CreateBuildInput) (*meroxa.Build, error)
	CreateSourceV2(ctx context.Context, input *meroxa.CreateSourceInputV2) (*meroxa.Source, error)
	GetBuild(ctx context.Context, uuid string) (*meroxa.Build, error)
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
	AddHeader(key, value string)
}

type environment struct {
	Name string
	UUID string
}

func (e *environment) nameOrUUID() string {
	switch {
	case e.UUID != "":
		return e.UUID
	case e.Name != "":
		return e.Name
	default:
		panic("bad state: name or uuid should be present")
	}
}

type Deploy struct {
	flags struct {
		Path                     string `long:"path" usage:"Path to the app directory (default is local directory)"`
		Environment              string `long:"env" usage:"environment (name or UUID) where application will be deployed to"`
		Spec                     string `long:"spec" usage:"Deployment specification version to use to build and deploy the app" hidden:"true"`
		SkipCollectionValidation bool   `long:"skip-collection-validation" usage:"Skips unique destination collection and looping validations"` //nolint:lll
		Verbose                  bool   `long:"verbose" usage:"Prints more logging messages" hidden:"true"`
	}

	client        deployApplicationClient
	config        config.Config
	logger        log.Logger
	turbineCLI    turbine.CLI
	appConfig     *turbine.AppConfig
	configAppName string
	appName       string
	gitBranch     string
	path          string
	lang          ir.Lang
	fnName        string
	specVersion   string
	env           *environment
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
	return "deploy [--path pwd]"
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

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}

// getAppSource will return the proper destination where the application source will be uploaded and fetched.
func (d *Deploy) getAppSource(ctx context.Context) (*meroxa.Source, error) {
	in := meroxa.CreateSourceInputV2{}
	if d.env != nil {
		in.Environment = &meroxa.EntityIdentifier{
			UUID: d.env.UUID,
			Name: d.env.Name,
		}
	}
	return d.client.CreateSourceV2(ctx, &in)
}

func (d *Deploy) getPlatformImage(ctx context.Context) (string, error) {
	d.logger.StartSpinner("\t", "Fetching Meroxa Platform source...")

	s, err := d.getAppSource(ctx)
	if err != nil {
		d.logger.Errorf(ctx, "\t 𐄂 Unable to fetch source")
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", err
	}
	d.logger.StopSpinnerWithStatus("Platform source fetched", log.Successful)

	err = d.UploadSource(ctx, s.PutUrl)
	if err != nil {
		return "", err
	}
	sourceBlob := meroxa.SourceBlob{Url: s.GetUrl}
	buildInput := &meroxa.CreateBuildInput{SourceBlob: sourceBlob}
	if d.env != nil {
		buildInput.Environment = &meroxa.EntityIdentifier{
			UUID: d.env.UUID,
			Name: d.env.Name,
		}
	}

	build, err := d.client.CreateBuild(ctx, buildInput)
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", err
	}
	d.logger.StartSpinner("\t", fmt.Sprintf("Building Meroxa Process image (%q)...", build.Uuid))

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
		time.Sleep(platformBuildPollDuration)
	}
}

func (d *Deploy) UploadSource(ctx context.Context, url string) error {
	var (
		err       error
		buildPath string
	)

	d.logger.StartSpinner("\t", fmt.Sprintf("Creating Dockerfile before uploading source in %s", d.path))
	buildPath, err = d.turbineCLI.CreateDockerfile(ctx, d.appName)
	if err != nil {
		return err
	}
	defer d.turbineCLI.CleanupDockerfile(d.logger, d.path)
	d.logger.StopSpinnerWithStatus("Dockerfile created", log.Successful)

	dFile := fmt.Sprintf("turbine-%s.tar.gz", d.appName)

	var buf bytes.Buffer
	if err = createTarAndZipFile(buildPath, &buf); err != nil {
		return err
	}

	d.logger.StartSpinner("\t", fmt.Sprintf("Creating %q in %q to upload to our build service...", buildPath, dFile))

	fileToWrite, err := os.OpenFile(dFile, os.O_CREATE|os.O_RDWR, os.FileMode(0o777)) //nolint:gomnd
	defer func(fileToWrite *os.File) {
		if err = fileToWrite.Close(); err != nil {
			panic(err.Error())
		}

		// remove .tar.gz file
		d.logger.StartSpinner("\t", fmt.Sprintf("Removing %q...", dFile))
		if err = os.Remove(dFile); err != nil {
			d.logger.StopSpinnerWithStatus(fmt.Sprintf("\t Something went wrong trying to remove %q", dFile), log.Failed)
		} else {
			d.logger.StopSpinnerWithStatus(fmt.Sprintf("Removed %q", dFile), log.Successful)
		}
	}(fileToWrite)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fileToWrite, &buf); err != nil {
		return err
	}
	d.logger.StopSpinnerWithStatus(fmt.Sprintf("%q successfully created in %q", dFile, buildPath), log.Successful)

	return uploadFile(ctx, d.logger, dFile, url)
}

// CreateTarAndZipFile creates a .tar.gz file from `src` on current directory.
func createTarAndZipFile(src string, buf io.Writer) error {
	// Grab the directory we care about (app's directory)
	appDir := filepath.Base(src)

	// Change to parent's app directory
	pwd, err := turbine.SwitchToAppDirectory(filepath.Dir(src))
	if err != nil {
		return err
	}

	zipWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(zipWriter)

	err = filepath.Walk(appDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if shouldSkipDir(fi) {
			return filepath.SkipDir
		}
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(file)
		if err := tarWriter.WriteHeader(header); err != nil { //nolint:govet
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		if !fi.IsDir() {
			var data *os.File
			data, err = os.Open(file)
			defer func(data *os.File) {
				err = data.Close()
				if err != nil {
					panic(err.Error())
				}
			}(data)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := tarWriter.Close(); err != nil {
		return err
	}
	if err := zipWriter.Close(); err != nil {
		return err
	}

	return os.Chdir(pwd)
}

func shouldSkipDir(fi os.FileInfo) bool {
	if !fi.IsDir() {
		return false
	}

	switch fi.Name() {
	case ".git", "fixtures", "node_modules":
		return true
	}

	return false
}

func uploadFile(ctx context.Context, logger log.Logger, filePath, url string) error {
	logger.StartSpinner("\t", "Uploading source...")

	var clientErr error
	var res *http.Response
	client := &http.Client{}

	retries := 3
	for retries > 0 {
		fh, err := os.Open(filePath)
		if err != nil {
			logger.StopSpinnerWithStatus("\t Failed to open source file", log.Failed)
			return err
		}
		defer func(fh *os.File) {
			fh.Close()
		}(fh)

		req, err := http.NewRequestWithContext(ctx, "PUT", url, fh)
		if err != nil {
			logger.StopSpinnerWithStatus("\t Failed to make new request", log.Failed)
			return err
		}

		fi, err := fh.Stat()
		if err != nil {
			logger.StopSpinnerWithStatus("\t Failed to stat source file", log.Failed)
			return err
		}
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Content-Type", "multipart/form-data")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Connection", "keep-alive")

		req.ContentLength = fi.Size()

		res, clientErr = client.Do(req)
		if clientErr != nil || (res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices) {
			retries--
		} else {
			break
		}
	}

	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	if clientErr != nil || (res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices) {
		logger.StopSpinnerWithStatus("\t Failed to upload build source file", log.Failed)
		if clientErr == nil {
			clientErr = fmt.Errorf("upload failed: %s", res.Status)
		}
		return clientErr
	}

	logger.StopSpinnerWithStatus("Source uploaded", log.Successful)
	return nil
}

func (d *Deploy) createDeployment(ctx context.Context, imageName, gitSha, specVersion string) (*meroxa.Deployment, error) {
	specStr, err := d.turbineCLI.GetDeploymentSpec(ctx, imageName)
	if err != nil {
		return nil, err
	}
	var spec map[string]interface{}
	if specStr != "" {
		if unmarshalErr := json.Unmarshal([]byte(specStr), &spec); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to parse deployment spec into json")
		}
	}
	if specVersion != "" {
		input := &meroxa.CreateDeploymentInput{
			Application: meroxa.EntityIdentifier{Name: d.appName},
			GitSha:      gitSha,
			SpecVersion: specVersion,
			Spec:        spec,
		}
		return d.client.CreateDeployment(ctx, input)
	}
	return nil, nil
}

// getAppImage will check what type of build needs to perform and ultimately will return the image name to use
// when deploying.
func (d *Deploy) getAppImage(ctx context.Context) (string, error) {
	d.logger.StartSpinner("\t", "Checking if application has processes...")

	needsToBuild, err := d.turbineCLI.NeedsToBuild(ctx)
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return "", err
	}

	// If no need to build, return empty imageName which won't be utilized by the deployment process anyway.
	if !needsToBuild {
		d.logger.StopSpinnerWithStatus("No need to create process image\n", log.Successful)
		return "", nil
	}

	d.logger.StopSpinnerWithStatus("Application processes found. Creating application image...", log.Successful)
	return d.getPlatformImage(ctx)
}

// validateLanguage stops execution of the deployment in case language is not supported.
// It also consolidates lang used in API in case user specified a supported language using an unexpected name.
func (d *Deploy) validateLanguage() error {
	switch d.lang {
	case "go", turbine.GoLang:
		d.lang = ir.GoLang
	case "js", turbine.JavaScript, turbine.NodeJs:
		d.lang = ir.JavaScript
	case "py", turbine.Python3, turbine.Python:
		d.lang = ir.Python
	case "rb", turbine.Ruby:
		d.lang = ir.Ruby
	default:
		return fmt.Errorf("language %q not supported. %s", d.lang, LanguageNotSupportedError)
	}
	return nil
}

// validateAppJSON will validate app JSON provided by customer has the right formation including supported language
// TODO: Extract some of this logic this turbine-core so we centralize language support in one place.
func (d *Deploy) validateAppJSON(ctx context.Context) error {
	var err error

	d.logger.Info(ctx, "Validating your app.json...")

	if d.path, err = turbine.GetPath(d.flags.Path); err != nil {
		return err
	}

	// exit early if app config is loaded
	if d.appConfig != nil {
		return d.validateLanguage()
	}

	d.lang, err = turbine.GetLangFromAppJSON(d.logger, d.path)
	if err != nil {
		return err
	}

	d.configAppName, err = turbine.GetAppNameFromAppJSON(d.logger, d.path)
	if err != nil {
		return err
	}
	if d.appConfig, err = turbine.ReadConfigFile(d.path); err != nil {
		return err
	}

	if d.gitBranch, err = turbine.GetGitBranch(d.path); err != nil {
		return err
	}
	d.appName = d.prepareAppName(ctx)

	return nil
}

func (d *Deploy) validateResources(ctx context.Context, rr []turbine.ApplicationResource) error {
	var errs []error
	validated := make(map[string]bool)

	wrapNotFound := func(name string, err error) error {
		if strings.Contains(err.Error(), "could not find") {
			return fmt.Errorf("could not find resource %q", name)
		}
		return err
	}

	for _, r := range rr {
		// dedup resources for validation
		if _, ok := validated[r.Name]; ok {
			continue
		}
		resource, err := d.client.GetResourceByNameOrID(ctx, r.Name)

		// order is important
		switch {
		case err != nil:
			errs = append(errs, wrapNotFound(r.Name, err))
		case resource.Status.State != meroxa.ResourceStateReady:
			errs = append(errs, fmt.Errorf("resource %q is not ready and usable", r.Name))
		// app is provisioned in common env, but resource was added at self hosted env
		case d.flags.Environment == "" && resource.Environment != nil:
			errs = append(errs, fmt.Errorf(
				"resource %q is in %q, but app is in common",
				r.Name,
				resource.Environment.Name,
			))
		// app is provisioned in an env, but resource is in common env
		case d.flags.Environment != "" && resource.Environment == nil:
			errs = append(errs, fmt.Errorf(
				"resource %q is not in app env %q, but in common",
				r.Name,
				d.flags.Environment,
			))
		// app is provisioned in an env, but resource is in different self hosted env
		case d.flags.Environment != "" && resource.Environment.Name != d.flags.Environment:
			errs = append(errs, fmt.Errorf(
				"resource %q is not in app env %q, but in %q",
				r.Name,
				d.flags.Environment,
				resource.Environment.Name,
			))
		}

		validated[r.Name] = true
	}

	return wrapErrors(errs)
}

func (d *Deploy) checkResourceAvailability(ctx context.Context) error {
	resourceCheckMessage := fmt.Sprintf("Checking resource availability for application %q (%s) before deployment...", d.appName, d.lang)

	d.logger.StartSpinner("\t", resourceCheckMessage)

	resources, err := d.turbineCLI.GetResources(ctx)
	if err != nil {
		return fmt.Errorf("unable to read resource definition from app: %s", err.Error())
	}

	if len(resources) == 0 {
		return errors.New("no resources defined in your Turbine app")
	}

	if err := d.validateResources(ctx, resources); err != nil {
		d.logger.StopSpinnerWithStatus("Resource availability check failed", log.Failed)
		return fmt.Errorf("%w;\n\n%s", err, resourceInvalidError)
	}

	if d.flags.SkipCollectionValidation {
		d.logger.StopSpinnerWithStatus("Can access your Turbine resources", log.Successful)
		return nil
	}

	if err := d.validateCollections(ctx, resources); err != nil {
		d.logger.StopSpinnerWithStatus("Resource availability check failed", log.Failed)
		return err
	}

	d.logger.StopSpinnerWithStatus("Can access your Turbine resources", log.Successful)
	return nil
}

func (d *Deploy) prepareDeployment(ctx context.Context) error {
	d.logger.Infof(ctx, "Preparing to deploy application %q...", d.appName)

	// check if resources exist and are ready
	err := d.checkResourceAvailability(ctx)
	if err != nil {
		return err
	}

	d.fnName, err = d.getAppImage(ctx)
	return err
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

// TODO: Eventually remove this and validate fast in Platform API.
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
			if r.Collection.Destination == "true" &&
				destinations[resourceCollectionPair{
					collectionName: r.Collection.Name,
					resourceName:   r.Name,
				}] {
				errMessage = fmt.Sprintf(
					"%s\n\tApplication resource %q with collection %q cannot be used as a destination. "+
						"It is also being used as a destination by another application %q.",
					errMessage,
					r.Name,
					r.Collection.Name,
					app.Name,
				)
			}
		}
	}

	return errMessage
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

func (d *Deploy) waitForDeployment(ctx context.Context, depUUID string) error {
	cctx, cancel := context.WithTimeout(ctx, minutesToWaitForDeployment*time.Minute)
	defer cancel()
	checkLogsMsg := "Check `meroxa apps logs` for further information"
	t := time.NewTicker(intervalCheckForDeployment)
	defer t.Stop()

	prevLine := ""

	for {
		select {
		case <-t.C:
			var deployment *meroxa.Deployment
			deployment, err := d.client.GetDeployment(ctx, d.appName, depUUID)
			if err != nil {
				return fmt.Errorf("couldn't fetch deployment status: %s", err.Error())
			}

			logs := strings.Split(deployment.Status.Details, "\n")

			if d.flags.Verbose {
				l := len(logs)
				if l > 0 && logs[l-1] != prevLine {
					prevLine = logs[l-1]
					d.logger.Info(ctx, "\t"+logs[l-1])
				}
			}

			switch {
			case deployment.Status.State == meroxa.DeploymentStateDeployed:
				return nil
			case deployment.Status.State == meroxa.DeploymentStateDeployingError:
				if !d.flags.Verbose {
					d.logger.Error(ctx, "\n")
					for _, l := range logs {
						d.logger.Errorf(ctx, "\t%s", l)
					}
				}
				return fmt.Errorf("\n %s", checkLogsMsg)
			}
		case <-cctx.Done():
			return fmt.Errorf(
				"your Turbine Application Deployment did not finish within %d minutes. %s",
				minutesToWaitForDeployment, checkLogsMsg)
		}
	}
}

// tearDownExistingResources will only delete the application and its associated entities if it hasn't been created
// or whether it's in a non-running state.
func (d *Deploy) tearDownExistingResources(ctx context.Context) error {
	app, _ := d.client.GetApplication(ctx, d.appName)

	if app != nil && app.Status.State != meroxa.ApplicationStateFailed {
		appExists := fmt.Sprintf("application %q exists in the %q state", d.appName, app.Status.State)
		msg := fmt.Sprintf(
			"%s\n\t. Use `meroxa apps remove %s` if you want to redeploy to this application",
			appExists,
			d.appName)
		return errors.New(msg)
	}
	resp, _ := d.client.DeleteApplicationEntities(ctx, d.appName)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	return nil
}

// validateEnvironmentFlagCompatibility takes care of some possible incompatibilities between flag values and languages.
// Turbine Ruby will use a spec'ed version internally so --spec is not required when using environments which is
// currently only available for those deployments using IR.
func (d *Deploy) validateEnvironmentFlagCompatibility() error {
	if d.flags.Spec == "" && d.env != nil {
		return fmt.Errorf(
			"please run `meroxa apps deploy` with `--spec %s` or `--spec %s` if you want to deploy to an environment",
			ir.SpecVersion_0_1_1, ir.SpecVersion_0_2_0)
	}
	return nil
}

// TODO: Once builds are done much faster we should move early checks like these to the Platform API.
func (d *Deploy) validateEnvExists(ctx context.Context) error {
	if _, err := d.client.GetEnvironment(ctx, d.env.nameOrUUID()); err != nil {
		if strings.Contains(err.Error(), "could not find environment") {
			return fmt.Errorf("environment %q does not exist", d.flags.Environment)
		}
		return fmt.Errorf("unable to retrieve environment %q: %w", d.flags.Environment, err)
	}
	return nil
}

// runValidations will perform the relevant validations in the order that's necessary.
func (d *Deploy) validate(ctx context.Context) error {
	if err := d.validateAppJSON(ctx); err != nil {
		return err
	}

	if err := d.validateEnvironmentFlagCompatibility(); err != nil {
		return err
	}

	if d.env != nil {
		return d.validateEnvExists(ctx)
	}

	return nil
}

//nolint:gocyclo,funlen
func (d *Deploy) Execute(ctx context.Context) error {
	var (
		app *meroxa.Application
		err error
	)

	if d.flags.Environment != "" {
		d.env = envFromFlag(d.flags.Environment)
	}

	if err = d.validate(ctx); err != nil {
		return err
	}

	if d.turbineCLI == nil {
		if d.turbineCLI, err = getTurbineCLIFromLanguage(d.logger, d.lang, d.path); err != nil {
			return err
		}
	}

	turbineLibVersion, err := d.turbineCLI.GetVersion(ctx)
	if err != nil {
		return err
	}
	addTurbineHeaders(d.client, d.lang, turbineLibVersion)

	if err = d.turbineCLI.GitChecks(ctx, d.logger, d.path); err != nil {
		return err
	}

	gitSha, err := d.turbineCLI.GetGitSha(ctx, d.path)
	if err != nil {
		return err
	}

	// SetupForDeploy will start the GRPC server for the specific language being deployed.
	cleanup, err := d.turbineCLI.SetupForDeploy(ctx, gitSha)
	if err != nil {
		return err
	}
	defer cleanup()

	// ⚠️ This is only until we re-deploy applications applying only the changes made
	if err = d.tearDownExistingResources(ctx); err != nil {
		return err
	}

	d.specVersion = d.flags.Spec

	// Make the latest version default for Ruby , Go , Python and JS
	if d.specVersion == "" {
		d.specVersion = ir.LatestSpecVersion
	}

	// Intermediate Representation Workflow
	if d.specVersion != "" {
		if err = ir.ValidateSpecVersion(d.specVersion); err != nil {
			return err
		}
		appInput := &meroxa.CreateApplicationInput{
			Name:     d.appName,
			Language: string(d.lang),
			GitSha:   gitSha,
		}
		// Creates application
		if d.env != nil {
			appInput.Environment = &meroxa.EntityIdentifier{
				Name: d.env.Name,
				UUID: d.env.UUID,
			}
		}
		app, err = d.client.CreateApplicationV2(ctx, appInput)

		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				msg := fmt.Sprintf("%s\n\tUse `meroxa apps remove %s` if you want to redeploy to this application", err, d.appName)
				return errors.New(msg)
			}
			return err
		}
	}

	err = d.prepareDeployment(ctx)
	if err != nil {
		return err
	}

	var deployment *meroxa.Deployment
	// Not deploying using IR
	deployMsg := fmt.Sprintf("Deploying application %q...", d.appName)

	// Hack: If not using IR as deployment type
	if d.specVersion == "" {
		d.logger.Infof(ctx, deployMsg)
	}

	if deployment, err = d.createDeployment(ctx, d.fnName, gitSha, d.specVersion); err != nil {
		return err
	}

	if d.specVersion == "" {
		// App was created by turbine libraries using pre-IR
		app, err = d.client.GetApplication(ctx, d.appName)
		if err != nil {
			return err
		}
	} else { // Deploying using IR
		// In verbose mode, we'll use spinners for each deployment step
		if d.flags.Verbose {
			d.logger.Info(ctx, deployMsg+"\n")
		} else {
			d.logger.StartSpinner("", deployMsg)
		}

		err = d.waitForDeployment(ctx, deployment.UUID)
		if err != nil {
			deploymentErroredMsg := "Couldn't complete the deployment"
			if !d.flags.Verbose {
				d.logger.StopSpinnerWithStatus(deploymentErroredMsg, log.Failed)
			} else {
				d.logger.Error(ctx, fmt.Sprintf("\t%s %s", d.logger.FailedMark(), deploymentErroredMsg))
			}
			return err
		}
	}

	dashboardURL := fmt.Sprintf("https://dashboard.meroxa.io/apps/%s/detail", d.appName)
	output := fmt.Sprintf("Application %q successfully deployed!\n\n  ✨ To view your application, visit %s",
		d.appName, dashboardURL)

	if d.flags.Verbose {
		d.logger.Info(ctx, fmt.Sprintf("\n\t%s %s", d.logger.SuccessfulCheck(), output))
	} else {
		d.logger.StopSpinnerWithStatus(output, log.Successful)
	}

	d.logger.JSON(ctx, app)
	return nil
}

func wrapErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	err := errs[0]
	for _, e := range errs[1:] {
		err = fmt.Errorf("%w; %v", err, e)
	}

	return err
}

// TODO: Reuse this everywhere it's using the environment flag. Maybe make this part of builder so it's automatic.
func envFromFlag(nameOrUUID string) *environment {
	_, err := uuid.Parse(nameOrUUID)
	if err != nil {
		return &environment{Name: nameOrUUID}
	}
	return &environment{UUID: nameOrUUID}
}
