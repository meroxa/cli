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
	"os"
	"os/exec"
	"path"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Run struct {
	flags struct {
		// `--lang` is not required unless language is not specified via app.json
		Lang string `long:"lang" short:"l" usage:"language to use (go | js)"`
		Path string `long:"path" usage:"path of application to run" required:"true"`
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
	Language string
}

func (*Run) Usage() string {
	return "run"
}

func (*Run) Docs() builder.Docs {
	return builder.Docs{
		Short: "Execute a Meroxa Data Application locally",
		Example: "meroxa apps run --path ../go-demo\n" +
			"meroxa apps run --path ../js-demo --lang js",
	}
}

func (r *Run) Flags() []builder.Flag {
	return builder.BuildFlags(&r.flags)
}

func (r *Run) buildGoApp(ctx context.Context) error {
	err := buildGoApp(ctx, r.logger, ".", false)
	if err != nil {
		return err
	}

	// apps name
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	projName := path.Base(pwd)

	cmd := exec.Command("./" + projName) //nolint:gosec
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	r.logger.Infof(ctx, "Running %s:", projName)
	r.logger.Info(ctx, string(stdout))

	return nil
}

func (r *Run) buildJSApp(ctx context.Context) error {
	// TODO: Handle this requirement https://github.com/meroxa/turbine-js.git being installed
	// (requirement being node)
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
	appConfig, err := r.readConfigFile()
	if err != nil {
		return err
	}

	lang := appConfig.Language

	if lang == "" {
		if r.flags.Lang == "" {
			return fmt.Errorf("flag --lang is required unless specified in your app.json")
		}
		lang = r.flags.Lang
	}

	switch lang {
	case "go", "golang":
		return r.buildGoApp(ctx)
	case "js", "javascript", "nodejs":
		return r.buildJSApp(ctx)
	default:
		return fmt.Errorf("language %q not supported. Currently, we support \"javascript\" and \"go\"", lang)
	}
}

func (r *Run) Logger(logger log.Logger) {
	r.logger = logger
}
