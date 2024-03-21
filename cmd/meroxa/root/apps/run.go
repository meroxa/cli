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

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Run struct {
	logger log.Logger

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
		Short: "Execute a Conduit Data Application locally",
		Long:  "meroxa apps run will build your app locally to then run it locally in --path.",
		Example: `meroxa apps run 			# assumes you run it from the app directory
meroxa apps run --path ../go-demo 	# it'll use lang defined in your app.json
`,
	}
}

func (r *Run) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Run) Flags() []builder.Flag {
	return builder.BuildFlags(&r.flags)
}

func (r *Run) Execute(_ context.Context) error {
	return nil
}
