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

package billing

import (
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/root/open"

	"github.com/spf13/cobra"
)

// TODO: Check how to disable parent flags (e.g.: --json)

// Cmd represents the `meroxa billing` command.
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "billing",
		Short: "Open your billing page in a web browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("meroxa open billing")
			err := open.Cmd().RunE(cmd, args)

			if err != nil {
				return err
			}

			return nil
		},
	}
}
