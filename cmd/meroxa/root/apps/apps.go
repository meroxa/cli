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

package apps

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/meroxa/cli/cmd/meroxa/builder"
)

type Apps struct{}

const (
	JavaScript = "javascript"
	GoLang     = "golang"
	NodeJs     = "nodejs"
)

var (
	_ builder.CommandWithDocs        = (*Apps)(nil)
	_ builder.CommandWithAliases     = (*Apps)(nil)
	_ builder.CommandWithSubCommands = (*Apps)(nil)
	_ builder.CommandWithHidden      = (*Apps)(nil)
	_ builder.CommandWithFeatureFlag = (*Apps)(nil)
)

func (*Apps) Aliases() []string {
	return []string{"app"}
}

func (*Apps) Usage() string {
	return "apps"
}

func (*Apps) Hidden() bool {
	return true
}

func (*Apps) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage Meroxa Data Applications",
	}
}

func (*Apps) FeatureFlag() (string, error) {
	return "applications", fmt.Errorf("no access to the Meroxa Data Applications feature")
}

func (*Apps) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Deploy{}),
		builder.BuildCobraCommand(&Init{}),
		builder.BuildCobraCommand(&Run{}),
		builder.BuildCobraCommand(&Create{}),
		builder.BuildCobraCommand(&Remove{}),
	}
}
