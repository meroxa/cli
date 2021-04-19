package global

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	Version string
	Config  *viper.Viper // TODO remove this global variable, read on demand
)

var (
	flagConfig  string
	flagDebug   bool
	FlagJSON    bool          // TODO make this private! do not use this variable from other packages
	FlagTimeout time.Duration // TODO make this private! do not use this variable from other packages
)

func RegisterGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&FlagJSON, "json", false, "output json")
	cmd.PersistentFlags().StringVar(&flagConfig, "config", "", "config file (default is $HOME/meroxa.env)")
	cmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "display any debugging information")
	cmd.PersistentFlags().DurationVar(&FlagTimeout, "timeout", time.Second*10, "set the client timeout")
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

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var err error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err != nil {
			// skip if we encountered an error along the way
			return
		}

		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --api-url to MEROXA_API_URL
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
			if err != nil {
				return
			}
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				return
			}
		}
	})
	return err
}
