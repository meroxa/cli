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

package cmd

import "github.com/spf13/cobra"

// RemoveCmd represents the `meroxa remove` command
func RemoveCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a component",
		Long: `Deprovision a component of the Meroxa platform, including pipelines,
 resources, and connectors`,
		SuggestFor: []string{"destroy", "delete"},
		Aliases:    []string{"rm", "delete"},
	}

	removeCmd.AddCommand(RemoveConnectorCmd())
	removeCmd.AddCommand(RemoveEndpointCmd())
	removeCmd.AddCommand(RemovePipelineCmd())
	removeCmd.AddCommand(RemoveResourceCmd())

	// This is to make acceptance tests happy
	var force bool
	removeCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "force delete without confirmation prompt")
	removeCmd.PersistentFlags().MarkHidden("force")

	return removeCmd
}
