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
	"github.com/meroxa/cli/cmd/meroxa/root/pipelines"
	"github.com/meroxa/cli/cmd/meroxa/root/resources"
	"github.com/spf13/cobra"
)

// updateCmd represents `meroxa update`.
func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a component",
		Long:  `Update a component of the Meroxa platform, including connectors`,
	}

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `[connector | endpoint | pipeline | resource] update` instead"
	}

	cmd.AddCommand(updateConnector())
	cmd.AddCommand(updatePipeline())
	cmd.AddCommand(updateResource())

	return cmd
}

// updateConnector represents `meroxa update connector` -> `meroxa connector update`.
func updateConnector() *cobra.Command {
	cmd := builder.BuildCobraCommand(&connectors.Update{})
	cmd.Use = "connector NAME --state pause | resume | restart"
	cmd.Short = "Update connector state"

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `connector update` instead"
	}

	return cmd
}

// updatePipeline represents `meroxa update pipeline` -> `meroxa pipeline update`.
func updatePipeline() *cobra.Command {
	cmd := builder.BuildCobraCommand(&pipelines.Update{})
	cmd.Use = "pipeline NAME"
	cmd.Short = "Update pipeline state"
	cmd.Example = "\n" +
		"meroxa update pipeline old-name --name new-name\n" +
		"meroxa update pipeline pipeline-name --state pause\n" +
		"meroxa update pipeline pipeline-name --metadata '{\"key\":\"value\"}'"

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `pipeline update` instead"
	}

	return cmd
}

// updateResource represents `meroxa update resource` -> `meroxa resource update`.
func updateResource() *cobra.Command {
	cmd := builder.BuildCobraCommand(&resources.Update{})
	cmd.Use = "resource NAME"
	cmd.Short = "Update a resource"
	cmd.Long = `Use the update command to update various Meroxa resources.`

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `resource update` instead"
	}

	return cmd
}
