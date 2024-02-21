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
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/v2/pkg/ir"
)

type Deploy struct {
	flags struct {
		Path string `long:"path" usage:"Path to the app directory (default is local directory)"`
		Spec string `long:"spec" usage:"Deployment specification version to use to build and deploy the app" hidden:"true"`
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

	// deployments needs to ensure enough time to complete
	if !global.ClientWithCustomTimeout() {
		d.client.SetTimeout(60 * time.Second)
	}
}

func (d *Deploy) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Deploy) getPlatformImage(ctx context.Context) error {
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

	d.appTarName = fmt.Sprintf("turbine-%s.tar.gz", d.appName)

	if err = d.buildAndRunImage(ctx); err != nil {
		return err
	}

	d.logger.StartSpinner("\t", fmt.Sprintf("Saving and uploading %q image", d.appName))

	if err = d.saveImageAsTar(ctx); err != nil {
		return err
	}

	d.logger.StopSpinnerWithStatus(fmt.Sprintf("%q successfully created in %q", d.appTarName, buildPath), log.Successful)
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

func (d *Deploy) buildAndRunImage(ctx context.Context) error {
	d.logger.StartSpinner("\t", fmt.Sprintf("Creating function image in %s", d.path))
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

func (d *Deploy) saveImageAsTar(ctx context.Context) error {
	app := fmt.Sprintf("%s:latest", d.appName)
	appTar := fmt.Sprintf("%s.tar", d.appName)
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"save",
		"-o",
		appTar,
		app,
	)

	cmd.Dir = d.path
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("\ndocker run failed: %v - %v", err.Error(), string(out))
	}

	err = d.gzipImageTar(ctx, appTar)
	return err
}

func (d *Deploy) gzipImageTar(_ context.Context, appTar string) error {
	reader, err := os.Open(filepath.Join(d.path, appTar))
	if err != nil {
		return err
	}

	filename := filepath.Base(filepath.Join(d.path, appTar))
	target := filepath.Join(d.path, d.appTarName)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join(d.path, appTar))
	return err
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

	input, err := d.createApplicationRequest(ctx)
	if err != nil {
		return err
	}

	/*	TODO: Enable when function integration is done

		d.logger.Infof(ctx, "Preparing to deploy application %q...", d.appName)

		if err = d.getPlatformImage(ctx); err != nil {
			return err
		}

		file := filepath.Join(d.path, d.appTarName)
		if _, err = os.Stat(file); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("turbine archive %q does not exist: %w", file, err)
			}

			return err
		}

		files := map[string]string{
			"imageArchive": file,
		}
	*/

	response, err := d.client.CollectionRequestMultipart(
		ctx,
		http.MethodPost,
		collectionName,
		"",
		input,
		nil,
		map[string]string{}, // TODO: change back to files from above
	)
	if err != nil {
		return err
	}

	apps := &Application{}
	err = json.NewDecoder(response.Body).Decode(&apps)
	if err != nil {
		return err
	}

	dashboardURL := fmt.Sprintf("%s/apps/%s/detail", global.GetMeroxaAPIURL(), apps.ID)
	output := fmt.Sprintf("Application %q successfully deployed!\n\n  ✨ To view your application, visit %s",
		d.appName, dashboardURL)

	d.logger.StopSpinnerWithStatus(output, log.Successful)

	return nil
}
