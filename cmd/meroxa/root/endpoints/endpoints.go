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

package endpoints

import (
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/spf13/cobra"
)

type Endpoints struct{}

var (
	_ builder.CommandWithAliases     = (*Endpoints)(nil)
	_ builder.CommandWithDocs        = (*Endpoints)(nil)
	_ builder.CommandWithSubCommands = (*Endpoints)(nil)
)

func (*Endpoints) Usage() string {
	return "endpoints"
}

func (*Endpoints) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage endpoints on Meroxa",
	}
}

func (*Endpoints) Aliases() []string {
	return []string{"endpoint"}
}

func (*Endpoints) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Create{}),
		builder.BuildCobraCommand(&Describe{}),
		builder.BuildCobraCommand(&List{}),
		builder.BuildCobraCommand(&Remove{}),
	}
}
