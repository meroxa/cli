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
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/spf13/cobra"
)

// LogoutCmd represents the `meroxa logout` command
func LogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "logout of the Meroxa platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: add confirmation
			global.Config.Set("ACCESS_TOKEN", "")
			global.Config.Set("REFRESH_TOKEN", "")
			err := global.Config.WriteConfig()
			if err != nil {
				return err
			}
			fmt.Println("Successfully logged out.")
			return nil
		},
	}
}
