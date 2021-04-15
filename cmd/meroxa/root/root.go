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
	"io/ioutil"
	"os"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/root/add"
	"github.com/meroxa/cli/cmd/meroxa/root/old"
	"github.com/meroxa/cli/log"
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
	var (
		cfgFile        string
		flagOutputJSON bool
		flagDebug      bool
	)

	cmd := &cobra.Command{
		Use:   "meroxa",
		Short: "The Meroxa CLI",
		Long: `The Meroxa CLI allows quick and easy access to the Meroxa data platform.

Using the CLI you are able to create and manage sophisticated data pipelines
with only a few simple commands. You can get started by listing the supported
resource types:

meroxa list resource-types`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			cfg, err := readConfig(cfgFile)
			if err != nil {
				return err
			}

			global.Config = cfg

			var (
				logLevel         = log.Info
				leveledLoggerOut = os.Stdout
				jsonLoggerOut    = ioutil.Discard
			)

			if flagOutputJSON {
				logLevel = log.Warn
				jsonLoggerOut = os.Stdout
			}
			if flagDebug {
				logLevel = log.Debug
			}

			global.Log = log.New(
				log.NewLeveledLogger(leveledLoggerOut, logLevel),
				log.NewJSONLogger(jsonLoggerOut),
			)

			return nil
		},

		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	cmd.PersistentFlags().BoolVar(&flagOutputJSON, "json", false, "output json")
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/meroxa.env)")
	cmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "display any debugging information")

	// Subcommands
	cmd.AddCommand(builder.BuildCobraCommand(&add.Add{}))
	cmd.AddCommand(old.ApiCmd())
	cmd.AddCommand(old.BillingCmd())
	cmd.AddCommand(old.CompletionCmd())
	cmd.AddCommand((&old.Connect{}).Command())
	cmd.AddCommand(old.CreateCmd())
	cmd.AddCommand(old.DescribeCmd())
	cmd.AddCommand(old.ListCmd())
	cmd.AddCommand(old.LoginCmd())
	cmd.AddCommand(old.LogoutCmd())
	cmd.AddCommand(old.LogsCmd())
	cmd.AddCommand(old.OpenCmd())
	cmd.AddCommand((&old.Remove{}).Command())
	cmd.AddCommand(old.UpdateCmd())
	cmd.AddCommand(old.VersionCmd())

	return cmd
}
