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

package deprecated

import "github.com/spf13/cobra"

// RegisterCommands Adds commands that follow an old CLI structure that's not longer in use.
func RegisterCommands(cmd *cobra.Command) {
	cmd.AddCommand(addCmd())
	cmd.AddCommand(createCmd())
	cmd.AddCommand(describeCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(logsCmd())
	cmd.AddCommand(removeCmd())
	cmd.AddCommand(updateCmd())
}
