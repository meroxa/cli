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
	"github.com/meroxa/cli/cmd/meroxa/root/endpoints"
	"github.com/meroxa/cli/cmd/meroxa/root/resources"
	"github.com/spf13/cobra"
)

// cmd represents the `meroxa describe` command.
func describeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a component",
		Long:  `Describe a component of the Meroxa data platform, including resources and connectors`,
	}

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `[connector | endpoint | pipeline | resource] describe` instead"
	}

	cmd.AddCommand(describeConnectorCmd())
	cmd.AddCommand(describeEndpointCmd())
	cmd.AddCommand(describeResourceCmd())

	return cmd
}

// describeConnectorCmd represents `meroxa describe connector` -> `meroxa connector describe`.
func describeConnectorCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&connectors.Describe{})
	cmd.Use = "connector NAME"
	cmd.Short = "Describe connector"

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `connector describe` instead"
	}

	return cmd
}

// describeEndpointCmd represents `meroxa describe endpoint` -> `meroxa endpoint describe`.
func describeEndpointCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&endpoints.Describe{})
	cmd.Use = "endpoint NAME"
	cmd.Short = "Describe endpoint"

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `endpoint describe` instead"
	}

	return cmd
}

// describeResourceCmd represents `meroxa describe resource` -> `meroxa endpoint describe`.
func describeResourceCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&resources.Describe{})
	cmd.Use = "resource NAME"
	cmd.Short = "Describe resource"

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `resource describe` instead"
	}

	return cmd
}
