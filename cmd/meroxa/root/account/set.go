/*
Copyright © 2022 Meroxa Inc

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

package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Set)(nil)
	_ builder.CommandWithArgs    = (*Set)(nil)
	_ builder.CommandWithLogger  = (*Set)(nil)
	_ builder.CommandWithClient  = (*Set)(nil)
	_ builder.CommandWithExecute = (*Set)(nil)
	_ builder.CommandWithConfig  = (*Set)(nil)
)

type Set struct {
	client listAccounts
	logger log.Logger
	config config.Config

	args struct {
		UUID string
	}
}

func (s *Set) Usage() string {
	return "set"
}

func (s *Set) Docs() builder.Docs {
	return builder.Docs{
		Short: "Set active project",
	}
}

func (s *Set) Client(client meroxa.Client) {
	s.client = client
}

func (s *Set) Config(cfg config.Config) {
	s.config = cfg
}

func (s *Set) Execute(ctx context.Context) error {
	accounts, err := s.client.ListAccounts(ctx)
	if err != nil {
		return err
	}

	found := false
	uuid := ""
	for _, account := range accounts {
		if s.args.UUID == account.Name ||
			s.args.UUID == account.UUID {
			found = true
			uuid = account.UUID
			break
		}
	}
	if !found {
		return fmt.Errorf("'%s' is an invalid project UUID", s.args.UUID)
	}
	s.config.Set(global.UserAccountUUID, uuid)

	return nil
}

func (s *Set) Logger(logger log.Logger) {
	s.logger = logger
}

func (s *Set) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires project UUID")
	}

	s.args.UUID = args[0]
	return nil
}
