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
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type RemoveEndpoint struct {
	name      string
	removeCmd *Remove
}

// RemoveEndpointClient represents the interface for meroxa client
type RemoveEndpointClient interface {
	DeleteEndpoint(ctx context.Context, name string) error
}

func (re *RemoveEndpoint) setArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires endpoint name\n\nUsage:\n  meroxa remove endpoint NAME [flags]")
	}
	// endpoint name
	re.name = args[0]

	re.removeCmd.componentType = "endpoint"
	re.removeCmd.confirmableName = re.name

	return nil
}

func (re *RemoveEndpoint) execute(ctx context.Context, c RemoveEndpointClient) error {
	return c.DeleteEndpoint(ctx, re.name)
}

func (re *RemoveEndpoint) output() {
	fmt.Printf("endpoint %s successfully removed\n", re.name)
}

// command represents the `meroxa remove endpoint` command
func (re *RemoveEndpoint) command() *cobra.Command {
	return &cobra.Command{
		Use:     "endpoint NAME",
		Aliases: []string{"endpoints"},
		Short:   "Remove endpoint",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return re.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()
			if err != nil {
				return err
			}

			// TODO: To be consistent with other commands, execute should return
			// also the component being removed
			err = re.execute(ctx, c)

			if err != nil {
				return err
			}

			// TODO: This usually send a component so it's shown in json format
			re.output()
			return nil
		},
	}
}
