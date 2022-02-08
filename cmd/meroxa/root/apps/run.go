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
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Run struct {
	flags struct {
		Lang string `long:"lang" short:"l" usage:"language to use (js | golang)" required:"true"`
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
		Example: "meroxa apps run ../go-demo --lang golang\n" +
			"meroxa apps run ../js-demo --lang js",
	}
}

func (r *Run) Flags() []builder.Flag {
	return builder.BuildFlags(&r.flags)
}

func (r *Run) Execute(ctx context.Context) error {
	var appPath string

	if p := r.flags.Path; p != "" {
		appPath = p
		if err := os.Chdir(appPath); err != nil {
			return err
		}
	} else {
		appPath = "."
	}

	appConfigPath := path.Join(appPath, "app.json")
	appConfigBytes, err := ioutil.ReadFile(appConfigPath)
	if err != nil {
		return err
	}
	var appConfig AppConfig
	if err := json.Unmarshal(appConfigBytes, &appConfig); err != nil {
		return err
	}

	if appConfig.Language == "go" {
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
	} else if appConfig.Language == "javascript" {
		cmd := exec.Command("npx", "turbine", "test")
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		r.logger.Info(ctx, string(stdout))

		return nil
	} else {
		return nil
	}
}

func (r *Run) Logger(logger log.Logger) {
	r.logger = logger
}
