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
	"os"
	"path/filepath"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
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
		Short: "Execute a Meroxa Data Application locally",
		Long: "meroxa apps run will build your app locally to then run it\n" +
			"locally on --path.",
		Example: "meroxa apps run # assumes you run it from the app directory\n" +
			"meroxa apps run --path ../js-demo --lang js # in case you didn't specify lang on your app.json" +
			"meroxa apps run --path ../go-demo # it'll use lang defined in your app.json",
	}
}

func (r *Run) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Run) Flags() []builder.Flag {
	return builder.BuildFlags(&r.flags)
}

func (r *Run) Execute(ctx context.Context) error {
	r.path = turbineCLI.GetPath(r.flags.Path)
	if r.path == "." {
		r.path, _ = filepath.Abs(r.path)
	} else if r.path == "" {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		r.path, _ = filepath.Abs(dir)
	}
	lang, err := turbineCLI.GetLang(r.flags.Lang, r.path)
	if err != nil {
		return err
	}

	switch lang {
	case "go", GoLang:
		return turbineCLI.RunGoApp(ctx, r.path, r.logger)
	case "js", JavaScript, NodeJs:
		return turbineCLI.BuildJSApp(ctx, r.logger)
	default:
		return fmt.Errorf("language %q not supported. Currently, we support \"javascript\" and \"go\"", lang)
	}
}
