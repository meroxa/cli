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

package transforms

import (
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/spf13/cobra"
)

type Transforms struct{}

var (
	_ builder.CommandWithAliases     = (*Transforms)(nil)
	_ builder.CommandWithDocs        = (*Transforms)(nil)
	_ builder.CommandWithSubCommands = (*Transforms)(nil)
)

func (*Transforms) Usage() string {
	return "transforms"
}

func (*Transforms) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage transforms on Meroxa",
	}
}

func (*Transforms) Aliases() []string {
	return []string{"transform"}
}

func (*Transforms) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&List{}),
	}
}
