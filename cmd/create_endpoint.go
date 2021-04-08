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
	"github.com/spf13/cobra"
)

var (
	flagEndpointCmdProtocol string
	flagEndpointCmdStream   string
)

// CreateEndpointCmd represents the `meroxa create endpoint` command
func CreateEndpointCmd() *cobra.Command {
	createEndpointCmd := &cobra.Command{
		Use:     "endpoint [NAME] [flags]",
		Aliases: []string{"endpoints"},
		Short:   "Create an endpoint",
		Long:    "Use create endpoint to expose an endpoint to a connector stream",
		Example: `
meroxa create endpoint my-endpoint --protocol http --stream my-stream`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client()
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), clientTimeOut)
			defer cancel()

			var name string
			if len(args) > 0 {
				name = args[0]
			}

			return c.CreateEndpoint(ctx, name, flagEndpointCmdProtocol, flagEndpointCmdStream)
		},
	}

	createEndpointCmd.Flags().StringVarP(&flagEndpointCmdProtocol, "protocol", "p", "", "protocol, value can be http or grpc (required)")
	createEndpointCmd.Flags().StringVarP(&flagEndpointCmdStream, "stream", "s", "", "stream name (required)")
	createEndpointCmd.MarkFlagRequired("protocol")
	createEndpointCmd.MarkFlagRequired("stream")
	return createEndpointCmd
}
