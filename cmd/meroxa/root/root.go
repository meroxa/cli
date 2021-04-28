/*
Copyright © 2020 Meroxa Inc

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

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/root/deprecated"
	"github.com/meroxa/cli/cmd/meroxa/root/deprecated/add"
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
	// TODO: Once we make a full transition to `subject-verb-object` remove these altogether
	cmd.AddCommand(builder.BuildCobraCommand(&add.Add{}))
	cmd.AddCommand(deprecated.CompletionCmd())
	cmd.AddCommand((&deprecated.Connect{}).Command())
	cmd.AddCommand(deprecated.CreateCmd())
	cmd.AddCommand(deprecated.DescribeCmd())
	cmd.AddCommand(deprecated.ListCmd())
	cmd.AddCommand(deprecated.LogsCmd())
	cmd.AddCommand(deprecated.OpenCmd())
	cmd.AddCommand((&deprecated.Remove{}).Command())
	cmd.AddCommand(deprecated.UpdateCmd())

	// v2
	cmd.AddCommand(APICmd())
	cmd.AddCommand(BillingCmd())
	cmd.AddCommand(LoginCmd())
	cmd.AddCommand(LogoutCmd())
	cmd.AddCommand(VersionCmd())
	cmd.AddCommand((&GetUser{}).Command()) // whoami

	return cmd
}
