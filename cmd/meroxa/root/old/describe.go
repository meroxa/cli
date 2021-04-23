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

// DescribeCmd represents the `meroxa describe` command.
func DescribeCmd() *cobra.Command {
	describeCmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a component",
		Long:  `Describe a component of the Meroxa data platform, including resources and connectors`,
	}

	describeCmd.AddCommand(DescribeResourceCmd())
	describeCmd.AddCommand(DescribeConnectorCmd())
	describeCmd.AddCommand(DescribeEndpointCmd())

	return describeCmd
}
