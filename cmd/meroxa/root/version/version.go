/*
Copyright © 2021 Meroxa Inc

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

package version

import (
	"context"
	"runtime"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
)

type Version struct {
	logger log.Logger
}

var (
	_ builder.CommandWithDocs    = (*Version)(nil)
	_ builder.CommandWithExecute = (*Version)(nil)
	_ builder.CommandWithLogger  = (*Version)(nil)
)

func (v *Version) Usage() string {
	return "version"
}

func (v *Version) Docs() builder.Docs {
	return builder.Docs{
		Short: "Display the Meroxa CLI version",
	}
}

func (v *Version) Logger(logger log.Logger) {
	v.logger = logger
}

func (v *Version) Execute(ctx context.Context) error {
	v.logger.Infof(ctx, "meroxa/%s %s/%s\n", global.Version, runtime.GOOS, runtime.GOARCH)
	return nil
}
