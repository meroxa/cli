/*
Copyright © 2021 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or impliee.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

var (
	_ builder.CommandWithDocs    = (*Describe)(nil)
	_ builder.CommandWithLogger  = (*Describe)(nil)
	_ builder.CommandWithExecute = (*Describe)(nil)
)

type Describe struct {
	logger log.Logger
}

func (e *Describe) Usage() string {
	return "describe"
}

func (e *Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Show Meroxa CLI configuration details",
		Example: "" +
			"$ meroxa config describe\n" +
			"Using meroxa config located in \"/Users/my-name/Library/Application Support/meroxa/config.env\n\n" +
			"access_token: c0f928b...c337a0d\n" +
			"actor: user@email.com\n" +
			"actor_uuid: c0f928ba-d40e-40c5-a7fa-cf281c337a0d\n" +
			"refresh_token: c337a0d...c0f928b\n",
		Long: "This command will return the content of your configuration file where you could find your `access_token` " +
			"and then `refresh_token` for your `meroxa login`. They're stored in the home directory on your machine. " +
			"On Unix, including MacOS, it's stored in `$HOME`, and on Windows is stored in `%USERPROFILE%`.",
	}
}

func (e *Describe) Execute(ctx context.Context) error {
	path := global.Config.ConfigFileUsed()

	var env struct {
		Path   string                 `json:"path"`
		Config map[string]interface{} `json:"config"`
	}

	env.Path = path
	env.Config = global.Config.AllSettings()

	cfgSettings := global.Config.AllSettings()
	cfgKeySettings := global.Config.AllKeys()

	sort.Strings(cfgKeySettings)

	e.logger.Infof(ctx, "Using meroxa config located in %q\n\n", path)

	for _, k := range cfgKeySettings {
		v := e.obfuscate(k, fmt.Sprintf("%s", cfgSettings[k]))
		e.logger.Infof(ctx, "%s: %s", k, v)
	}

	e.logger.JSON(ctx, env)
	return nil
}

func (e *Describe) Logger(logger log.Logger) {
	e.logger = logger
}

func (e *Describe) obfuscate(key, value string) string {
	if !strings.Contains(key, "token") {
		return value
	}

	if len(value) < 5 { // nolint:gomnd
		// hide whole text
		return strings.Repeat("*", len(value))
	}

	const (
		maxVisibleLen = 7
	)

	visibleLen := (len(value) - 3) / 2
	if visibleLen > maxVisibleLen {
		visibleLen = maxVisibleLen
	}

	return value[:visibleLen] + "..." + value[len(value)-visibleLen:]
}
