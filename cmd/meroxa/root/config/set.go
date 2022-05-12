/*
Copyright © 2022 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or impliee.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"regexp"
	"strings"
)

var (
	_ builder.CommandWithDocs    = (*Set)(nil)
	_ builder.CommandWithLogger  = (*Set)(nil)
	_ builder.CommandWithExecute = (*Set)(nil)
	_ builder.CommandWithArgs    = (*Set)(nil)
	_ builder.CommandWithConfig  = (*Set)(nil)
)

type Set struct {
	logger log.Logger
	config config.Config
	args   struct {
		keys map[string]string
	}
}

func (s *Set) Usage() string {
	return "set"
}

func (s *Set) Docs() builder.Docs {
	return builder.Docs{
		Short: "Update your Meroxa CLI configuration file with a specific key=value",
		Example: "" +
			"$ meroxa config set DisableUpdateNotification=true\n" +
			"$ meroxa config set DISABLE_UPDATE_NOTIFICATION=true\n" +
			"$ meroxa config set OneKey=true AnotherKey=false\n" +
			"$ meroxa config set ApiUrl=https://staging.meroxa.com\n",
		Long: "This command will let you update your Meroxa configuration file to customize your CLI experience." +
			"You can check the presence of this file by running `meroxa config describe`, or even provide your own using `--config my-other-cfg-file`" +
			"A key with a format such as MyKey will be converted automatically to as MY_KEY.",
	}
}

func (s *Set) Execute(ctx context.Context) error {
	for k, v := range s.args.keys {
		s.logger.Infof(ctx, "Updating your Meroxa configuration file with %s=%s...", k, v)
		s.config.Set(k, v)

	}
	s.logger.Info(ctx, "Done!")
	// run over each key and write them to a file
	// show confirmation message
	return nil
}

func (s *Set) Logger(logger log.Logger) {
	s.logger = logger
}

func (s *Set) Config(cfg config.Config) {
	s.config = cfg
}

// TODO: Implement (and Test!)
func (s *Set) normalizeKey(key string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(key, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}

func (s *Set) validateAndAssignKeyValue(kv string) error {
	nkv := strings.Split(kv, "=")

	if len(nkv) != 2 {
		return fmt.Errorf("a key=value needs to contain at least and only one `=` sign")
	}

	k := s.normalizeKey(nkv[0])
	v := nkv[1]

	s.args.keys[k] = v
	return nil
}

func (s *Set) ParseArgs(args []string) error {
	s.args.keys = make(map[string]string)

	var err error
	if len(args) < 1 {
		return errors.New("requires at least one KEY=VALUE pair (example: meroxa config set KEY=VALUE)")
	}

	for _, a := range args {
		err = s.validateAndAssignKeyValue(a)
		if err != nil {
			return err
		}
	}
	return nil
}
