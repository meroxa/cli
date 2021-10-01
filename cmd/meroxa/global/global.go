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

package global

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	Version string
	Config  *viper.Viper
)

var (
	flagCLIConfigFile string
	flagAPIURL        string
	flagDebug         bool
	flagTimeout       time.Duration
	flagJSON          bool
)

const (
	AccessTokenEnv       = "ACCESS_TOKEN"
	ActorEnv             = "ACTOR"
	ActorUUIDEnv         = "ACTOR_UUID"
	CasedDebugEnv        = "CASED_DEBUG"
	CasedPublishKeyEnv   = "CASED_PUBLISH_KEY"
	PublishMetricsEnv    = "PUBLISH_METRICS"
	RefreshTokenEnv      = "REFRESH_TOKEN"
	UserFeatureFlagsEnv  = "USER_FEATURE_FLAGS"
	UserInfoUpdatedAtEnv = "USER_INFO_UPDATED_AT"
)

func RegisterGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "output json")

	cmd.PersistentFlags().StringVar(&flagCLIConfigFile, "cli-config-file", "", "meroxa configuration file")
	cmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", "", "API url")
	cmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "display any debugging information")
	cmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", time.Second*10, "set the duration of the client timeout in seconds") // nolint:gomnd,lll

	if err := cmd.PersistentFlags().MarkHidden("api-url"); err != nil {
		panic(fmt.Sprintf("could not mark flag as hidden: %v", err))
	}
}

func PersistentPreRunE(cmd *cobra.Command) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}
	Config = cfg

	err = bindFlags(cmd, Config)
	if err != nil {
		return err
	}

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable).
func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var err error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err != nil {
			// skip if we encountered an error along the way
			return
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.GetString(f.Name)
			err = cmd.Flags().Set(f.Name, val)
			if err != nil {
				return
			}
		}
	})
	return err
}
