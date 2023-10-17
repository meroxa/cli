/*
Copyright © 2022 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or impliew.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"context"
	"os"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
)

type WhoAmI struct {
	logger log.Logger
	config config.Config
}

var (
	_ builder.CommandWithDocs    = (*WhoAmI)(nil)
	_ builder.CommandWithExecute = (*WhoAmI)(nil)
	_ builder.CommandWithLogger  = (*WhoAmI)(nil)
)

func (w *WhoAmI) Usage() string {
	return "whoami"
}

func (w *WhoAmI) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Display the current logged in user\n",
		Example: "meroxa whoami",
	}
}

func (w *WhoAmI) Logger(logger log.Logger) {
	w.logger = logger
}

func (w *WhoAmI) Execute(ctx context.Context) error {
	email := os.Getenv(global.TenantEmailAddress)
	w.logger.Infof(ctx, "%s", email)
	w.logger.JSON(ctx, email)

	return nil
}
