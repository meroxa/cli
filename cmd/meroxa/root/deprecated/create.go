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

// ⚠️ WARN: ⚠️
//
// The following commands will be removed once we decide to stop adding support for commands that don't follow
// the `subject-verb-object` design.
package deprecated

import (
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/root/connectors"
	"github.com/meroxa/cli/cmd/meroxa/root/endpoints"
	"github.com/meroxa/cli/cmd/meroxa/root/pipelines"
	"github.com/spf13/cobra"
)

// createCmd represents the `meroxa create` command.
func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Meroxa pipeline components",
		Long: `Use the create command to create various Meroxa pipeline components
including connectors.`,
	}

	if global.DeprecateV1Commands() {
		cmd.Deprecated = "use `[connector | endpoint | pipeline | resource] create` instead"
	}

	cmd.AddCommand(createConnectorCmd())
	cmd.AddCommand(createEndpoint())
	cmd.AddCommand(createPipelineCmd())

	return cmd
}

// createConnectorCmd represents `meroxa create connector` -> `meroxa connector create`.
func createConnectorCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&connectors.Create{})

	cmd.Use = "connector [NAME] [flags]"
	cmd.Short = "Create a connector"
	cmd.Long = "Use create connector to create a connector from a source (--from) or to a destination (--to)"
	cmd.Example = "\n" +
		"meroxa create connector [NAME] --from pg2kafka --input accounts \n" +
		"meroxa create connector [NAME] --to pg2redshift --input orders # --input will be the desired stream \n" +
		"meroxa create connector [NAME] --to pg2redshift --input orders --pipeline my-pipeline\n"

	cmd.Aliases = []string{"connectors"}

	if global.DeprecateV1Commands() {
		cmd.Deprecated = "use `connector create` instead"
	}

	return cmd
}

// createEndpoint represents `meroxa create endpoint` -> `meroxa endpoint create`.
func createEndpoint() *cobra.Command {
	cmd := builder.BuildCobraCommand(&endpoints.Create{})
	cmd.Use = "endpoint [NAME] [flags]"
	cmd.Short = "Create an endpoint"
	cmd.Long = "Use create endpoint to expose an endpoint to a connector stream"
	cmd.Example = `
meroxa create endpoint my-endpoint --protocol http --stream my-stream`
	cmd.Aliases = []string{"endpoints"}

	if global.DeprecateV1Commands() {
		cmd.Deprecated = "use `endpoint create` instead"
	}

	return cmd
}

// createPipelineCmd represents `meroxa create pipeline` -> `meroxa pipelines create`.
func createPipelineCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&pipelines.Create{})

	cmd.Use = "pipeline NAME" //nolint:goconst
	cmd.Short = "Create a pipeline"
	cmd.Aliases = []string{"pipelines"}

	if global.DeprecateV1Commands() {
		cmd.Deprecated = "use `pipelines create` instead"
	}

	return cmd
}
