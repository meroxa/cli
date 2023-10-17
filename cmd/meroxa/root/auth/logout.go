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

package auth

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
)

var (
	_ builder.CommandWithDocs    = (*Logout)(nil)
	_ builder.CommandWithExecute = (*Logout)(nil)
	_ builder.CommandWithLogger  = (*Logout)(nil)
	_ builder.CommandWithConfig  = (*Logout)(nil)
)

type Logout struct {
	logger log.Logger
	config config.Config
}

func (l *Logout) Usage() string {
	return "logout"
}

func (l *Logout) Docs() builder.Docs {
	return builder.Docs{
		Short: "Clears local login credentials of the Meroxa Platform",
	}
}

func (l *Logout) Logger(logger log.Logger) {
	l.logger = logger
}

func (l *Logout) Config(cfg config.Config) {
	l.config = cfg
}

func (l *Logout) Execute(ctx context.Context) error {
	l.config.Set(global.AccessTokenEnv, "")
	l.config.Set(global.RefreshTokenEnv, "")
	l.config.Set(global.UserFeatureFlagsEnv, "")

	l.logger.Infof(ctx, "Successfully logged out.")
	return nil
}
