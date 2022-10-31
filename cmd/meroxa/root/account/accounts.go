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
	"github.com/spf13/cobra"

	"github.com/meroxa/cli/cmd/meroxa/builder"
)

type Account struct{}

var (
	_ builder.CommandWithDocs        = (*Account)(nil)
	_ builder.CommandWithAliases     = (*Account)(nil)
	_ builder.CommandWithSubCommands = (*Account)(nil)
)

func (*Account) Aliases() []string {
	return []string{"accounts"}
}

func (*Account) Usage() string {
	return "account"
}

func (*Account) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage Meroxa Accounts",
	}
}

func (*Account) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&List{}),
		builder.BuildCobraCommand(&Set{}),
	}
}
