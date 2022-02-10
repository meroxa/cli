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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Run struct {
	flags struct {
		// `--lang` is not required unless language is not specified via app.json
		Lang string `long:"lang" short:"l" usage:"language to use (go | js)"`
		Path string `long:"path" usage:"path of application to run"`
	}

	logger log.Logger
}

var (
	_ builder.CommandWithDocs    = (*Run)(nil)
	_ builder.CommandWithFlags   = (*Run)(nil)
	_ builder.CommandWithExecute = (*Run)(nil)
	_ builder.CommandWithLogger  = (*Run)(nil)
)

type AppConfig struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

func (*Run) Usage() string {
	return "run"
}

func (*Run) Docs() builder.Docs {
	return builder.Docs{
		Short: "Execute a Meroxa Data Application locally",
		Long: "meroxa apps run will build your app locally to then run it\n" +
			"locally based on --path.",
		Example: "meroxa apps run # assumes you run it from the app directory\n" +
			"meroxa apps run --path ../js-demo --lang js # in case you didn't specify lang on your app.json" +
			"meroxa apps run --path ../go-demo # it'll use lang defined in your app.json",
	}
}

func (r *Run) Flags() []builder.Flag {
	return builder.BuildFlags(&r.flags)
}

func (r *Run) buildGoApp(ctx context.Context, appPath string) error {
	// grab current location to use it as project name
	appName := path.Base(appPath)

	// building is a requirement prior to running for go apps
	err := buildGoApp(ctx, r.logger, appPath, appName, false)
	if err != nil {
		return err
	}

	cmd := exec.Command("./" + appName) //nolint:gosec
	cmd.Dir = appPath
	fmt.Println("COMMAND: ", cmd.Dir, cmd)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return err
	}
	r.logger.Infof(ctx, "Running app %q:", appName)
	r.logger.Info(ctx, string(stdout))

	return nil
}

func (r *Run) buildJSApp(ctx context.Context) error {
	// TODO: Handle the requirement of https://github.com/meroxa/turbine-js.git being installed
	// cd into the path first
	cmd := exec.Command("npx", "turbine", "test")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	r.logger.Info(ctx, string(stdout))
	return nil
}

func (r *Run) readConfigFile() (AppConfig, error) {
	var appConfig AppConfig

	appPath := r.flags.Path
	appConfigPath := path.Join(appPath, "app.json")
	appConfigBytes, err := ioutil.ReadFile(appConfigPath)
	if err != nil {
		return appConfig, fmt.Errorf("%v\n"+
			"Applications to run require an app.json file\n"+
			"Check out this example: https://github.com/meroxa/valve/blob/ah/demo2/examples/simple/app.json", err)
	}
	if err := json.Unmarshal(appConfigBytes, &appConfig); err != nil {
		return appConfig, err
	}

	return appConfig, nil
}

func (r *Run) Execute(ctx context.Context) error {
	var appPath string

	switch {
	case r.flags.Path != "":
		appPath = r.flags.Path
	default:
		appPath = "."
	}

	appConfig, err := r.readConfigFile()
	if err != nil {
		return err
	}

	lang := appConfig.Language

	if lang == "" {
		if r.flags.Lang == "" {
			return fmt.Errorf("flag --lang is required unless lang is specified in your app.json")
		}
		lang = r.flags.Lang
	}

	switch lang {
	case "go", "golang":
		return r.buildGoApp(ctx, appPath)
	case "js", "javascript", "nodejs":
		return r.buildJSApp(ctx)
	default:
		return fmt.Errorf("language %q not supported. Currently, we support \"javascript\" and \"go\"", lang)
	}
}

func (r *Run) Logger(logger log.Logger) {
	r.logger = logger
}
