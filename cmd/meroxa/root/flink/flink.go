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

package flink

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/meroxa/cli/cmd/meroxa/builder"
)

type Job struct{}

var (
	_ builder.CommandWithDocs        = (*Job)(nil)
	_ builder.CommandWithSubCommands = (*Job)(nil)
	_ builder.CommandWithFeatureFlag = (*Job)(nil)
	_ builder.CommandWithHidden      = (*Job)(nil)
)

func (*Job) Usage() string {
	return "flink"
}

func (*Job) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage Flink Jobs",
	}
}

func (*Job) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Deploy{}),
		builder.BuildCobraCommand(&Remove{}),
		builder.BuildCobraCommand(&List{}),
	}
}

func (*Job) FeatureFlag() (string, error) {
	return "flink", fmt.Errorf(`no access to the Meroxa Flink Jobs feature`)
}

func (*Job) Hidden() bool {
	return true
}
