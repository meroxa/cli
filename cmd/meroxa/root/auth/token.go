/*
Copyright © 2021 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://wwt.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or impliet.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"context"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/global"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Token struct {
	logger log.Logger

	flags struct {
		Path bool `long:"path" usage:"returns path where token is stored"`
	}
}

var (
	_ builder.CommandWithDocs    = (*Token)(nil)
	_ builder.CommandWithLogger  = (*Token)(nil)
	_ builder.CommandWithExecute = (*Token)(nil)
	_ builder.CommandWithFlags   = (*Token)(nil)
)

func (t *Token) Usage() string {
	return "token"
}

func (t *Token) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Display the current logged in token\n",
		Example: "meroxa auth token",
	}
}

func (t *Token) Logger(logger log.Logger) {
	t.logger = logger
}

func (t *Token) Flags() []builder.Flag {
	return builder.BuildFlags(&t.flags)
}

func (t *Token) Execute(ctx context.Context) error {
	type UserToken struct {
		path  string
		token string
	}

	cfg, err := global.ReadConfig()

	if err != nil {
		return err
	}

	if t.flags.Path {
		t.logger.Infof(ctx, cfg.ConfigFileUsed())
	} else {
		token := cfg.Get("ACCESS_TOKEN")
		t.logger.Infof(ctx, fmt.Sprintf("%s", token))
	}

	return nil
}
