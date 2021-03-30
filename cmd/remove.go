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
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

type Remove struct {
	confirmableName string
	componentType   string
	force           bool
}

// confirmRemoved will prompt for confirmation
func (r *Remove) confirmRemove(stdin io.Reader, val string) error {
	reader := bufio.NewReader(stdin)
	fmt.Printf("To proceed, type %s or re-run this command with --force\n▸ ", val)
	input, _ := reader.ReadString('\n')

	if val != strings.TrimSuffix(input, "\n") {
		if r.componentType != "" {
			return errors.New(fmt.Sprintf("removing %s not confirmed", r.componentType))
		} else {
			return errors.New("removing value not confirmed")
		}
	}
	return nil
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

	cmd.AddCommand(RemoveEndpointCmd())
	cmd.AddCommand(RemovePipelineCmd())

	rc := &RemoveConnector{removeCmd: r}
	cmd.AddCommand(rc.command())
	rr := &RemoveResource{removeCmd: r}
	cmd.AddCommand(rr.command())

	// Make sure all subcommands will have a confirmation prompt or make use of --force
	for _, c := range cmd.Commands() {
		r.addConfirmation(c)
	}

	r.setFlags(cmd)
	return cmd
}

func (r *Remove) addConfirmation(subCmd *cobra.Command) {
	preRunE := subCmd.PreRunE

	subCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		err := preRunE(cmd, args)
		if err != nil {
			return err
		}

		// print and confirm
		if !flagRootOutputJSON {
			fmt.Printf("Removing %s...\n", r.confirmableName)
		}

		// prompts for confirmation when --force is not set
		if !r.force {
			return r.confirmRemove(os.Stdin, r.confirmableName)
		}

		return nil
	}
}
