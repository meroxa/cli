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

<<<<<<< HEAD
<<<<<<< HEAD:cmd/meroxa/root/billing/billing.go
package billing
=======
package open
>>>>>>> 48bdb79 (refactor: Commands to their own pkgs):cmd/meroxa/root/open/open.go

import (
	"github.com/meroxa/cli/cmd/meroxa/root/open"
)

<<<<<<< HEAD:cmd/meroxa/root/billing/billing.go
type Billing struct {
	open.Billing
=======
// Cmd represents the `meroxa open` command.
func Cmd() *cobra.Command {
	openCmd := &cobra.Command{
		Use:   "open",
		Short: "Open in a web browser",
	}

	openCmd.AddCommand(SubCmd())
	return openCmd
>>>>>>> 48bdb79 (refactor: Commands to their own pkgs):cmd/meroxa/root/open/open.go
=======
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
>>>>>>> 48bdb79 (refactor: Commands to their own pkgs)
}
