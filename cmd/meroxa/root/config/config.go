/*
Copyright © 2022 Meroxa Inc

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
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/spf13/cobra"
)

type Config struct {
	Describe
}

var (
	_ builder.CommandWithAliases     = (*Config)(nil)
	_ builder.CommandWithDocs        = (*Config)(nil)
	_ builder.CommandWithSubCommands = (*Config)(nil)
)

func (c *Config) Usage() string {
	return "config"
}

func (*Config) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage your Meroxa CLI configuration",
	}
}

func (*Config) Aliases() []string {
	return []string{"cfg"}
}

func (*Config) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Describe{}),
		builder.BuildCobraCommand(&Set{}),
	}
}
