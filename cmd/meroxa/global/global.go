package global

import (
	"github.com/meroxa/cli/log"
	"github.com/spf13/viper"
)

var (
	Version string
	Log     log.Logger
	Config  *viper.Viper
)
