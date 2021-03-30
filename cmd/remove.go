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
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

type Remove struct {
	force bool
}

// confirmRemoved will prompt for confirmation or will check the `--force` flag value
func (r *Remove) confirmRemove(stdin io.Reader, val string) bool {
	if !r.force {
		reader := bufio.NewReader(stdin)
		fmt.Printf("To proceed, type %s or re-run this command with --force\n▸ ", val)
		input, _ := reader.ReadString('\n')
		return val == strings.TrimSuffix(input, "\n")
	}

	return r.force
}

func (r *Remove) setFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&r.force, "force", "f", false, "force delete without confirmation prompt")
}

// RemoveCmd represents the `meroxa remove` command
func (r *Remove) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a component",
		Long: `Deprovision a component of the Meroxa platform, including pipelines,
 resources, and connectors`,
		SuggestFor: []string{"destroy", "delete"},
		Aliases:    []string{"rm", "delete"},
	}

	cmd.AddCommand(RemoveConnectorCmd())
	cmd.AddCommand(RemoveEndpointCmd())
	cmd.AddCommand(RemovePipelineCmd())

	rr := RemoveResource{removeCmd: r}
	cmd.AddCommand(rr.command())

	r.setFlags(cmd)
	return cmd
}
