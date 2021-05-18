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

// removeCmd represents `meroxa remove`.
func removeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a component",
		Long: `Deprovision a component of the Meroxa platform, including pipelines,
 resources, and connectors`,
		SuggestFor: []string{"destroy", "delete"},
		Aliases:    []string{"rm", "delete"},
	}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `[connector | endpoint | pipeline | resource] remove` instead"
	}

	cmd.AddCommand(removeConnectorCmd())
	cmd.AddCommand(removeEndpointCmd())
	cmd.AddCommand(removePipelineCmd())
	cmd.AddCommand(removeResourceCmd())
	return cmd
}

// removeConnectorCmd represents `meroxa remove connector` -> `meroxa connector remove`.
func removeConnectorCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&connectors.Remove{})
	cmd.Use = "connector NAME"
	cmd.Short = "Remove connector"
	cmd.Aliases = []string{"connectors"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `connector remove` instead"
	}

	return cmd
}

// removeEndpointCmd represents `meroxa remove endpoint` -> `meroxa endpoint remove`.
func removeEndpointCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&endpoints.Remove{})
	cmd.Use = "endpoint NAME"
	cmd.Short = "Remove endpoint"
	cmd.Aliases = []string{"endpoints"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `endpoint remove` instead"
	}

	return cmd
}

// removePipelineCmd represents `meroxa remove pipeline` -> `meroxa pipeline remove`.
func removePipelineCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&pipelines.Remove{})
	cmd.Use = "pipeline NAME"
	cmd.Short = "Remove pipeline"
	cmd.Aliases = []string{"pipelines"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `pipeline remove` instead"
	}

	return cmd
}

// removeResourceCmd represents `meroxa remove resource` -> `meroxa resource remove`.
func removeResourceCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&resources.Remove{})
	cmd.Use = "resource NAME"
	cmd.Short = "Remove resource"
	cmd.Aliases = []string{"resources"}

	if global.ShowDeprecationWarning() {
		cmd.Deprecated = "use `resource remove` instead"
	}

	return cmd
}
