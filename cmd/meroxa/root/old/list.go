/*
Copyright © 2020 Meroxa Inc

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

package old

import (
	"github.com/spf13/cobra"
)

// ListCmd represents the `meroxa list` command.
func ListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List components",
		Long: `List the components of the Meroxa platform, including pipelines,
 resources, connectors, etc... You may also filter by type.`,
	}

	listCmd.AddCommand((&ListConnectors{}).command())
	listCmd.AddCommand(ListEndpointsCmd())
	listCmd.AddCommand(ListResourceTypesCmd())
	listCmd.AddCommand(ListPipelinesCmd())
	listCmd.AddCommand(ListResourcesCmd())
	listCmd.AddCommand(ListTransformsCmd())

	return listCmd
}
