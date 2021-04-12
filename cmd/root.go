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
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	// The name of our config file, without the file extension because viper supports many different config file languages.
	defaultConfigFilename = "meroxa"

	// The environment variable prefix of all environment variables bound to our command line flags.
	envPrefix     = "MEROXA"
	clientTimeOut = 10 * time.Second
)

var (
	meroxaVersion                 string
	cfgFile                       string
	source                        string
	destination                   string
	flagRootOutputJSON, flagDebug bool
	meroxaCmd                     *cobra.Command
	cfg                           *viper.Viper
)

// RootCmd represents the base command when called without any subcommands
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "meroxa",
		Short: "The Meroxa CLI",
		Long: `The Meroxa CLI allows quick and easy access to the Meroxa data platform.

Using the CLI you are able to create and manage sophisticated data pipelines
with only a few simple commands. You can get started by listing the supported
resource types:

meroxa list resource-types`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initConfig(cmd)
		},
		TraverseChildren: true,
	}

	rootCmd.PersistentFlags().BoolVar(&flagRootOutputJSON, "json", false, "output json")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/meroxa.env)")
	rootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "display any debugging information")

	rootCmd.SilenceUsage = true
	rootCmd.DisableAutoGenTag = true

	// Subcommands
	rootCmd.AddCommand(AddCmd())
	rootCmd.AddCommand(ApiCmd())
	rootCmd.AddCommand(BillingCmd())
	rootCmd.AddCommand(CompletionCmd())
	rootCmd.AddCommand((&Connect{}).command())
	rootCmd.AddCommand(CreateCmd())
	rootCmd.AddCommand(DescribeCmd())
	rootCmd.AddCommand(ListCmd())
	rootCmd.AddCommand(LoginCmd())
	rootCmd.AddCommand(LogoutCmd())
	rootCmd.AddCommand(LogsCmd())
	rootCmd.AddCommand(OpenCmd())
	rootCmd.AddCommand((&Remove{}).command())
	rootCmd.AddCommand(UpdateCmd())
	rootCmd.AddCommand(VersionCmd())

	return rootCmd
}

// Execute fills in the version based on main.version by goreleaser.
// This is called by main.main()
func Execute(version string) {
	meroxaVersion = version

	if err := meroxaCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	meroxaCmd = RootCmd()
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	cfg = viper.New()

	if cfgFile != "" {
		// Use config file from the flag.
		cfg.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Set the base name of the config file, without the file extension.
		cfg.SetConfigName(defaultConfigFilename)
		cfg.AddConfigPath(home)
	}
	cfg.SetConfigType("env")
	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := cfg.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable MEROXA_NUMBER. This helps
	// avoid conflicts.
	cfg.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	cfg.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, cfg)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --api-url to MEROXA_API_URL
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
