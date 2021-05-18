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
	"github.com/meroxa/cli/cmd/meroxa/root/pipelines"
	"github.com/meroxa/cli/cmd/meroxa/root/resources"
	"github.com/spf13/cobra"
)

// listCmd represents the `meroxa list` command.
func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List components",
		Long: `List the components of the Meroxa platform, including pipelines,
 resources, connectors, etc... You may also filter by type.`,
	}

	cmd.AddCommand(listConnectorsCmd())
	cmd.AddCommand(listEndpointsCmd())
	cmd.AddCommand(listPipelinesCmd())
	cmd.AddCommand(listResourcesCmd())
	cmd.AddCommand(listResourceTypesCmd()) // WIP
	cmd.AddCommand(listTransformsCmd())

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `[connectors | endpoints | pipelines | resources | transforms] list` instead"
	}

	return cmd
}

// listConnectorsCmd represents `meroxa list connectors` -> `meroxa connectors list`.
func listConnectorsCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&connectors.List{})
	cmd.Use = "connectors"
	cmd.Short = "List connectors"
	cmd.Aliases = []string{"connector"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `connectors list` instead"
	}

	return cmd
}

// listEndpointsCmd represents `meroxa list endpoints` -> `meroxa endpoints list`.
func listEndpointsCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&endpoints.List{})
	cmd.Use = "endpoints"
	cmd.Short = "List endpoints"
	cmd.Aliases = []string{"endpoint"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `endpoints list` instead"
	}

	return cmd
}

// listPipelinesCmd represents `meroxa list pipelines` -> `meroxa pipelines list`.
func listPipelinesCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&pipelines.List{})
	cmd.Use = "pipelines"
	cmd.Short = "List pipelines"
	cmd.Aliases = []string{"pipeline"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `pipelines list` instead"
	}

	return cmd
}

// listResourcesCmd represents `meroxa list resources` -> `meroxa resources list`.
func listResourcesCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&resources.List{})
	cmd.Use = "resources"
	cmd.Short = "List resources"
	cmd.Aliases = []string{"resource"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `resources list` instead"
	}

	return cmd
}

// listResourceTypesCmd represents `meroxa list resource-types` -> `meroxa resources list`.
func listResourceTypesCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&resources.List{ListTypes: true})
	cmd.Use = "resource-types"
	cmd.Short = "List resource-types"
	cmd.Aliases = []string{"resource-type"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `resources list --types` instead"
	}

	return cmd
}

// listTransformsCmd represents `meroxa list transforms` -> `meroxa transforms list`.
func listTransformsCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&resources.List{})
	cmd.Use = "transforms"
	cmd.Short = "List transforms"
	cmd.Aliases = []string{"transform"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `transforms list` instead"
	}

	return cmd
}
