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

import (
	"github.com/spf13/cobra"
)

var (
	con            string // connector name
	res            string // resource name
	cfgString      string
	metadataString string
	input          string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Meroxa pipeline components",
	Long: `Use the create command to create various Meroxa pipeline components
including connectors.`,
}


func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.AddCommand(createConnectorCmd)
	createCmd.AddCommand(createPipelineCmd)
	createCmd.AddCommand(createEndpointCmd)
}

