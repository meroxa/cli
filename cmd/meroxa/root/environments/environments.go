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

package environments

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Environments struct {
	logger log.Logger
}

var (
	_ builder.CommandWithAliases     = (*Environments)(nil)
	_ builder.CommandWithDocs        = (*Environments)(nil)
	_ builder.CommandWithFeatureFlag = (*Environments)(nil)
	_ builder.CommandWithSubCommands = (*Environments)(nil)
)

func (*Environments) Usage() string {
	return "environments"
}

func (*Environments) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage environments on Meroxa",
	}
}

func (*Environments) Aliases() []string {
	return []string{"env", "environment"}
}

func (*Environments) FeatureFlag() (string, error) {
	return "environments", fmt.Errorf(`no access to the Meroxa self-hosted environments feature.
Sign up for the Beta here: https://share.hsforms.com/1Uq6UYoL8Q6eV5QzSiyIQkAc2sme`)
}

func (e *Environments) Logger(logger log.Logger) {
	e.logger = logger
}

func (*Environments) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Create{}),
		builder.BuildCobraCommand(&Describe{}),
		builder.BuildCobraCommand(&List{}),
		builder.BuildCobraCommand(&Remove{}),
		builder.BuildCobraCommand(&Update{}),
		builder.BuildCobraCommand(&Repair{}),
	}
}
