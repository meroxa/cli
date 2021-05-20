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

/* ⚠️ WARN ⚠️

The following commands will be removed once we decide to stop adding support for commands that don't follow
the `subject-verb-object` design.

*/

package deprecated

import (
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/root/connectors"
	"github.com/spf13/cobra"
)

// logsCmd represents the `meroxa logs` command.
func logsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Print logs for a component",
	}

	cmd.AddCommand(logsConnectorCmd())

	return cmd
}

// logsConnectorCmd represents `meroxa logs connector` -> `meroxa connector logs`.
func logsConnectorCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&connectors.Logs{})
	cmd.Use = "connector NAME"
	cmd.Short = "Print logs for a connector"

	if global.DeprecateV1Commands() {
		cmd.Deprecated = "use `connector logs` instead"
	}

	return cmd
}
