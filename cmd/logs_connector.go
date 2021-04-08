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
	"io"
	"os"

	"github.com/spf13/cobra"
)

// LogsConnectorCmd represents the `meroxa logs connector` command
func LogsConnectorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connector NAME",
		Short: "Print logs for a connector",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires connector name\n\nUsage:\n  meroxa logs connector NAME")
			}
			connector := args[0]

			c, err := client()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			resp, err := c.GetConnectorLogs(ctx, connector)
			if err != nil {
				return err
			}

			_, err = io.Copy(os.Stdout, resp.Body)
			if err != nil {
				return err
			}

			os.Stdout.Write([]byte("\n"))

			return nil
		},
	}
}
