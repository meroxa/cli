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
	"context"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/spf13/cobra"
)

type CreateEndpointClient interface {
	CreateEndpoint(ctx context.Context, name, protocol, stream string) error
}

type CreateEndpoint struct {
	name, protocol, stream string
}

func (ce *CreateEndpoint) setArgs(args []string) error {
	if len(args) > 0 {
		ce.name = args[0]
	}

	return nil
}

func (ce *CreateEndpoint) setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ce.protocol, "protocol", "p", "", "protocol, value can be http or grpc (required)")
	cmd.Flags().StringVarP(&ce.stream, "stream", "s", "", "stream name (required)")
	cmd.MarkFlagRequired("protocol")
	cmd.MarkFlagRequired("stream")
}

func (ce *CreateEndpoint) output() {
	fmt.Println("Endpoint successfully created!")
}

func (ce *CreateEndpoint) execute(ctx context.Context, c CreateEndpointClient) error {
	fmt.Println("Creating endpoint...")
	return c.CreateEndpoint(ctx, ce.name, ce.protocol, ce.stream)
}

// CreateEndpointCmd represents the `meroxa create endpoint` command
func (ce *CreateEndpoint) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "endpoint [NAME] [flags]",
		Aliases: []string{"endpoints"},
		Short:   "Create an endpoint",
		Long:    "Use create endpoint to expose an endpoint to a connector stream",
		Example: `
meroxa create endpoint my-endpoint --protocol http --stream my-stream`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ce.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := global.NewClient()
			if err != nil {
				return err
			}

			err = ce.execute(cmd.Context(), c)

			if err != nil {
				return err
			}

			ce.output()

			return nil
		},
	}

	ce.setFlags(cmd)
	return cmd
}
