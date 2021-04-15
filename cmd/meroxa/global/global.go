package global

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version string
	Config  *viper.Viper // TODO remove this global variable, read on demand
)

var (
	flagConfig string
	flagDebug  bool
	FlagJSON   bool // TODO make this private! do not use this variable from other packages
)

func RegisterGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&FlagJSON, "json", false, "output json")
	cmd.PersistentFlags().StringVar(&flagConfig, "config", "", "config file (default is $HOME/meroxa.env)")
	cmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "display any debugging information")
}
