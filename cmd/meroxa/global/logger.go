package global

import (
	"io/ioutil"
	"os"

	"github.com/meroxa/cli/log"
)

func NewLogger() log.Logger {
	var (
		logLevel         = log.Info
		leveledLoggerOut = os.Stdout
		jsonLoggerOut    = ioutil.Discard
	)

	if FlagJSON {
		logLevel = log.Warn
		jsonLoggerOut = os.Stdout
	}
	if flagDebug {
		logLevel = log.Debug
	}

	return log.New(
		log.NewLeveledLogger(leveledLoggerOut, logLevel),
		log.NewJSONLogger(jsonLoggerOut),
	)
}
