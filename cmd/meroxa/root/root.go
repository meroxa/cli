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

package root

import (
	"context"
	"os"

	"github.com/meroxa/cli/cmd/meroxa/root/api"
	"github.com/meroxa/cli/cmd/meroxa/root/auth"
	"github.com/meroxa/cli/cmd/meroxa/root/billing"
	"github.com/meroxa/cli/cmd/meroxa/root/connectors"
	"github.com/meroxa/cli/cmd/meroxa/root/endpoints"
	"github.com/meroxa/cli/cmd/meroxa/root/env"
	"github.com/meroxa/cli/cmd/meroxa/root/open"
	"github.com/meroxa/cli/cmd/meroxa/root/pipelines"
	"github.com/meroxa/cli/cmd/meroxa/root/resources"
	"github.com/meroxa/cli/cmd/meroxa/root/transforms"
	"github.com/meroxa/cli/cmd/meroxa/root/version"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/spf13/cobra"
)

func Run() {
	ctx := context.Background()

	rootCmd := Cmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

// Cmd represents the base command when called without any subcommands.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meroxa",
		Short: "The Meroxa CLI",
		Long: `The Meroxa CLI allows quick and easy access to the Meroxa Data Platform.

Using the CLI you are able to create and manage sophisticated data pipelines
with only a few simple commands. You can get started by listing the supported
resource types:

meroxa resources list --types
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return global.PersistentPreRunE(cmd)
		},
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		TraverseChildren:  true,
	}

	global.RegisterGlobalFlags(cmd)

	// Subcommands
	cmd.AddCommand(CompletionCmd()) // Coming from Cobra

	cmd.AddCommand(builder.BuildCobraCommand(&api.API{}))
	cmd.AddCommand(builder.BuildCobraCommand(&auth.Auth{}))
	cmd.AddCommand(builder.BuildCobraCommand(&billing.Billing{}))
	cmd.AddCommand(builder.BuildCobraCommand(&connectors.Connect{}))
	cmd.AddCommand(builder.BuildCobraCommand(&connectors.Connectors{}))
	cmd.AddCommand(builder.BuildCobraCommand(&endpoints.Endpoints{}))
	cmd.AddCommand(builder.BuildCobraCommand(&env.Env{}))
	cmd.AddCommand(builder.BuildCobraCommand(&open.Open{}))
	cmd.AddCommand(builder.BuildCobraCommand(&pipelines.Pipelines{}))
	cmd.AddCommand(builder.BuildCobraCommand(&resources.Resources{}))
	cmd.AddCommand(builder.BuildCobraCommand(&transforms.Transforms{}))
	cmd.AddCommand(builder.BuildCobraCommand(&version.Version{}))

	setAliases(cmd)

	return cmd
}

// setAliases includes command to root but not shown in help
// e.g.: `meroxa login` -> `meroxa auth login`.
func setAliases(cmd *cobra.Command) {
	aliases := map[string]builder.Command{
		"login":  &auth.Login{},
		"logout": &auth.Logout{},
		"whoami": &auth.WhoAmI{},
	}

	for v, c := range aliases {
		cc := builder.BuildCobraCommand(c)
		cc.Hidden = true
		cc.Use = v
		cmd.AddCommand(cc)
	}
}
