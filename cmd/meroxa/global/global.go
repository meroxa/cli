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
	flagConfig  string
	flagAPIURL  string
	flagDebug   bool
	flagTimeout time.Duration
	flagJSON    bool
)

func DeprecateV1Commands() bool {
	return true
}

func RegisterGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "output json")
	cmd.PersistentFlags().StringVar(&flagConfig, "config", "", "config file")
	cmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", "", "API url")
	cmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "display any debugging information")
	cmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", time.Second*10, "set the client timeout") // nolint:gomnd

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
