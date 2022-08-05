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
	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine_cli/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine_cli/javascript"
	turbinepy "github.com/meroxa/cli/cmd/meroxa/turbine_cli/python"
	"github.com/meroxa/cli/log"
)

type Run struct {
	flags struct {
		// `--lang` is not required unless language is not specified via app.json
		Lang string `long:"lang" short:"l" usage:"language to use (go | js)"`
		Path string `long:"path" usage:"path of application to run"`
	}

	path   string
	logger log.Logger
}

var (
	_ builder.CommandWithDocs    = (*Run)(nil)
	_ builder.CommandWithFlags   = (*Run)(nil)
	_ builder.CommandWithExecute = (*Run)(nil)
	_ builder.CommandWithLogger  = (*Run)(nil)
)

func (*Run) Usage() string {
	return "run"
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
	r.path, err = turbineCLI.GetPath(r.flags.Path)
	if err != nil {
		return err
	}
	lang, err := turbineCLI.GetLang(ctx, r.logger, r.flags.Lang, r.path)
	if err != nil {
		return err
	}

	switch lang {
	case GoLang:
		return turbineGo.Run(ctx, r.path, r.logger)
	case "js", JavaScript, NodeJs:
		return turbineJS.Build(ctx, r.logger, r.path)
	case "py", Python3, Python:
		return turbinepy.Run(ctx, r.logger, r.path)
	default:
		return fmt.Errorf("language %q not supported. %s", lang, LanguageNotSupportedError)
	}
}
