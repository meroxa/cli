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
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine/javascript"
	turbinePy "github.com/meroxa/cli/cmd/meroxa/turbine/python"
	turbineRB "github.com/meroxa/cli/cmd/meroxa/turbine/ruby"
	"github.com/meroxa/cli/log"
)

type Run struct {
	path   string
	config *turbine.AppConfig

	logger     log.Logger
	turbineCLI turbine.CLI

	flags struct {
		Path string `long:"path" usage:"path of application to run"`
	}
}

var (
	_ builder.CommandWithDocs    = (*Run)(nil)
	_ builder.CommandWithFlags   = (*Run)(nil)
	_ builder.CommandWithExecute = (*Run)(nil)
	_ builder.CommandWithLogger  = (*Run)(nil)
)

func (*Run) Usage() string {
	return "run [--path pwd]"
}

func (*Run) Docs() builder.Docs {
	return builder.Docs{
		Short: "Execute a Turbine Data Application locally",
		Long:  "meroxa apps run will build your app locally to then run it locally in --path.",
		Example: `meroxa apps run 			# assumes you run it from the app directory
meroxa apps run --path ../go-demo 	# it'll use lang defined in your app.json
`,
		Beta: true,
	}
}

func (r *Run) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Run) Flags() []builder.Flag {
	return builder.BuildFlags(&r.flags)
}

func (r *Run) Execute(ctx context.Context) error {
	var err error
	if r.config == nil {
		r.path, err = turbine.GetPath(r.flags.Path)
		if err != nil {
			return err
		}
		var config turbine.AppConfig
		config, err = turbine.ReadConfigFile(r.path)
		r.config = &config
		if err != nil {
			return err
		}
	}

	switch lang := r.config.Language; lang {
	case "go", turbine.GoLang:
		if r.turbineCLI == nil {
			r.turbineCLI = turbineGo.New(r.logger, r.path)
		}
		err = r.turbineCLI.Run(ctx)
		turbineGo.RunCleanup(ctx, r.logger, r.path, r.config.Name)
		return err
	case "js", turbine.JavaScript, turbine.NodeJs:
		if r.turbineCLI == nil {
			r.turbineCLI = turbineJS.New(r.logger, r.path)
		}
		return r.turbineCLI.Run(ctx)
	case "py", turbine.Python3, turbine.Python:
		if r.turbineCLI == nil {
			r.turbineCLI = turbinePy.New(r.logger, r.path)
		}
		return r.turbineCLI.Run(ctx)
	case "rb", turbine.Ruby:
		if r.turbineCLI == nil {
			r.turbineCLI = turbineRB.New(r.logger, r.path)
		}
		return r.turbineCLI.Run(ctx)
	default:
		return fmt.Errorf("language %q not supported. %s", lang, LanguageNotSupportedError)
	}
}
