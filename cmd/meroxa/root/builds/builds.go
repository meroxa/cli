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

package builds

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/meroxa/cli/cmd/meroxa/builder"
)

var (
	_ builder.CommandWithDocs        = (*Builds)(nil)
	_ builder.CommandWithAliases     = (*Builds)(nil)
	_ builder.CommandWithHidden      = (*Builds)(nil)
	_ builder.CommandWithFeatureFlag = (*Builds)(nil)
	_ builder.CommandWithSubCommands = (*Builds)(nil)
)

type Builds struct{}

func (*Builds) Usage() string {
	return "builds"
}

func (*Builds) Aliases() []string {
	return []string{"build"}
}

func (*Builds) Hidden() bool {
	return true
}

func (*Builds) FeatureFlag() (string, error) {
	return "turbine", fmt.Errorf("no access to the Meroxa Data Application feature")
}

func (*Builds) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage Process Builds on Meroxa",
	}
}

func (*Builds) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Describe{}),
		builder.BuildCobraCommand(&Logs{}),
	}
}
