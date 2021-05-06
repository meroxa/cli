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

	"github.com/meroxa/cli/cmd/meroxa/root/pipelines"
	"github.com/meroxa/cli/cmd/meroxa/root/resources"

	"github.com/meroxa/cli/cmd/meroxa/root/endpoints"

	"github.com/meroxa/cli/cmd/meroxa/root/connectors"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/root/deprecated"
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
		Long: `The Meroxa CLI allows quick and easy access to the Meroxa data platform.

Using the CLI you are able to create and manage sophisticated data pipelines
with only a few simple commands. You can get started by listing the supported
resource types:

meroxa list resource-types`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			deprecated.FlagRootOutputJSON = global.FlagJSON
			return global.PersistentPreRunE(cmd)
		},
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		TraverseChildren:  true,
	}

	global.RegisterGlobalFlags(cmd)

	// Subcommands

	// v1
	if v, ok := os.LookupEnv("MEROXA_V2"); !ok || ok && v != "only" {
		// TODO: Once we make a full transition to `subject-verb-object` remove the `deprecated` pkg altogether
		deprecated.RegisterCommands(cmd)
	}

	// v2
	cmd.AddCommand(APICmd())
	cmd.AddCommand(BillingCmd())
	cmd.AddCommand(builder.BuildCobraCommand(&connectors.Connect{}))
	cmd.AddCommand(LoginCmd())
	cmd.AddCommand(LogoutCmd())
	cmd.AddCommand(VersionCmd())
	cmd.AddCommand((&GetUser{}).Command()) // whoami

	// New commands following `subject-verb-object` only shown if using `MEROXA_V2`)
	if _, ok := os.LookupEnv("MEROXA_V2"); ok {
		cmd.AddCommand(builder.BuildCobraCommand(&connectors.Connectors{}))
		cmd.AddCommand(builder.BuildCobraCommand(&endpoints.Endpoints{}))
		cmd.AddCommand(builder.BuildCobraCommand(&pipelines.Pipelines{}))
		cmd.AddCommand(builder.BuildCobraCommand(&resources.Resources{}))
	}

	return cmd
}
