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

	if flagJSON {
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
