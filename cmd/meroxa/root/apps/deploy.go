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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/ir"
)

type Deploy struct {
	flags struct {
		Path                     string `long:"path" usage:"Path to the app directory (default is local directory)"`
		Spec                     string `long:"spec" usage:"Deployment specification version to use to build and deploy the app" hidden:"true"`
		SkipCollectionValidation bool   `long:"skip-collection-validation" usage:"Skips unique destination collection and looping validations"` //nolint:lll
		Verbose                  bool   `long:"verbose" usage:"Prints more logging messages" hidden:"true"`
	}

	client        global.BasicClient
	config        config.Config
	logger        log.Logger
	turbineCLI    turbine.CLI
	appConfig     *turbine.AppConfig
	configAppName string
	appName       string
	gitBranch     string
	path          string
	lang          ir.Lang
	fnName        string // is this still necessary?
	appTarName    string
	specVersion   string
	gitSha        string
}

var (
	_ builder.CommandWithBasicClient = (*Deploy)(nil)
	_ builder.CommandWithConfig      = (*Deploy)(nil)
	_ builder.CommandWithDocs        = (*Deploy)(nil)
	_ builder.CommandWithExecute     = (*Deploy)(nil)
	_ builder.CommandWithFlags       = (*Deploy)(nil)
	_ builder.CommandWithLogger      = (*Deploy)(nil)
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

func (d *Deploy) BasicClient(client global.BasicClient) {
	d.client = client
}

func (d *Deploy) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Deploy) createAndRemoveZipFile(buildPath string) error {
	dFile := fmt.Sprintf("turbine-%s.tar.gz", d.appName)

	var buf bytes.Buffer
	if err := createTarAndZipFile(buildPath, &buf); err != nil {
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
	d.appTarName = dFile
	d.fnName = dFile
	d.logger.StopSpinnerWithStatus(fmt.Sprintf("%q successfully created in %q", dFile, buildPath), log.Successful)
	return nil
}

func (d *Deploy) runDockerImage(ctx context.Context) error {
	// docker run -d --rm -p 8080:80 -t simple-with-process-mdpx-demo
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"run",
		"-d",
		"--rm",
		"-p",
		"8080:8080",
		"-t",
		d.appName,
	)

	cmd.Dir = d.path
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("\ndocker run failed: %v", string(out))
	}
	return nil
}

func (d *Deploy) buildAndRunImageForMDPX(ctx context.Context) error {
	err := d.mvnPkgCmd(ctx)
	if err != nil {
		return fmt.Errorf("failed building maven: %w", err)
	}

	d.logger.StartSpinner("\t", fmt.Sprintf("Creating function image in %s", d.path))
	_, err = d.turbineCLI.CreateDockerfile(ctx, "")
	if err != nil {
		return fmt.Errorf("failed creating Dockerfile: %w", err)
	}
	defer d.turbineCLI.CleanupDockerfile(d.logger, d.path)

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"build",
		"-t",
		d.appName,
		".",
	)

	cmd.Dir = d.path
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("\ndocker build failed: %v", string(out))
	}

	d.fnName = d.appName

	return d.runDockerImage(ctx)
}

// TODO: Update this before shipping to MDPX.
func (d *Deploy) getPlatformImage(ctx context.Context) error {
	var (
		err error
	)

	d.logger.StartSpinner("\t", fmt.Sprintf("Creating Dockerfile before uploading source in %s", d.path))
	buildPath, err := d.turbineCLI.CreateDockerfile(ctx, d.appName)
	if err != nil {
		return err
	}
	defer d.turbineCLI.CleanupDockerfile(d.logger, d.path)
	d.logger.StopSpinnerWithStatus("Dockerfile created", log.Successful)

	return d.createAndRemoveZipFile(buildPath)
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

// getAppImage will check what type of build needs to perform and ultimately will return the image name to use
// when deploying.
func (d *Deploy) getAppImage(ctx context.Context) error {
	d.logger.StartSpinner("\t", "Checking if application has processes...")

	needsToBuild, err := d.turbineCLI.NeedsToBuild(ctx)
	if err != nil {
		d.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	// If no need to build, return empty imageName which won't be utilized by the deployment process anyway.
	if !needsToBuild {
		d.logger.StopSpinnerWithStatus("No need to create process image\n", log.Successful)
		return nil
	}

	d.logger.StopSpinnerWithStatus("Application processes found. Creating function image...", log.Successful)

	// Classic Platform
	//return d.getPlatformImage(ctx)

	return d.buildAndRunImageForMDPX(ctx)
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
	case turbine.Java:
		d.lang = ir.Java
	default:
		return newLangUnsupportedError(d.lang)
	}
	return nil
}

// readFromAppJSON will validate app JSON provided by customer has the right formation including supported language.
func (d *Deploy) readFromAppJSON(ctx context.Context) error {
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

func (d *Deploy) prepareDeployment(ctx context.Context) error {
	d.logger.Infof(ctx, "Preparing to deploy application %q...", d.appName)
	return d.getAppImage(ctx)
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

func (d *Deploy) assignDeploymentValues(ctx context.Context) error {
	var err error

	// Always default to the latest spec version.
	d.specVersion = ir.LatestSpecVersion

	if d.flags.Spec != "" {
		d.specVersion = d.flags.Spec
	}

	if err = ir.ValidateSpecVersion(d.specVersion); err != nil {
		return err
	}

	if err = d.readFromAppJSON(ctx); err != nil {
		return err
	}

	if d.turbineCLI, err = getTurbineCLIFromLanguage(d.logger, d.lang, d.path); err != nil {
		return err
	}

	return nil
}

func (d *Deploy) getGitInfo(ctx context.Context) error {
	var err error

	if err = d.turbineCLI.CheckUncommittedChanges(ctx, d.logger, d.path); err != nil {
		return err
	}

	d.gitSha, err = d.turbineCLI.GetGitSha(ctx, d.path)
	return err
}

func (d *Deploy) createApplicationRequest(ctx context.Context) (*Application, error) {
	specStr, err := d.turbineCLI.GetDeploymentSpec(ctx, d.fnName)
	if err != nil {
		return nil, err
	}
	var spec map[string]interface{}
	if specStr != "" {
		if unmarshalErr := json.Unmarshal([]byte(specStr), &spec); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to parse deployment spec into json")
		}
	}
	return &Application{
		Name:        d.appName,
		SpecVersion: d.specVersion,
		Spec:        spec,
	}, nil
}

func (d *Deploy) mvnPkgCmd(ctx context.Context) error {
	cmd := exec.CommandContext(
		ctx,
		"mvn",
		"clean",
		"package",
	)

	cmd.Dir = d.path
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mvn failed: %v", string(out))
	}

	return nil
}

func (d *Deploy) Execute(ctx context.Context) error {
	var err error

	if err = d.assignDeploymentValues(ctx); err != nil {
		return err
	}

	turbineLibVersion, err := d.turbineCLI.GetVersion(ctx)
	if err != nil {
		return err
	}
	addTurbineHeaders(d.client, d.lang, turbineLibVersion)

	if err = d.getGitInfo(ctx); err != nil { //nolint:shadow
		return err
	}

	gracefulStop, err := d.turbineCLI.StartGrpcServer(ctx, d.gitSha)
	if err != nil {
		return err
	}
	defer gracefulStop()

	if err = d.prepareDeployment(ctx); err != nil {
		return err
	}

	input, err := d.createApplicationRequest(ctx)
	if err != nil {
		return err
	}
	if _, err = d.client.CollectionRequest(ctx, "POST", collectionName, "", input, nil); err != nil {
		return err
	}

	dashboardURL := fmt.Sprintf("https://dashboard.meroxa.io/apps/%s/detail", d.appName)
	output := fmt.Sprintf("Application %q successfully deployed!\n\n  ✨ To view your application, visit %s",
		d.appName, dashboardURL)

	if d.flags.Verbose {
		d.logger.Info(ctx, fmt.Sprintf("\n\t%s %s", d.logger.SuccessfulCheck(), output))
	} else {
		d.logger.StopSpinnerWithStatus(output, log.Successful)
	}

	return nil
}
