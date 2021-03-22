package cmd

import (
	"context"
	"fmt"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

func DescribeEndpointCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "endpoint <name>",
		Aliases: []string{"endpoints"},
		Short:   "Describe endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires endpoint name\n\nUsage:\n  meroxa describe endpoint <name> [flags]")
			}
			name := args[0]

			c, err := client()
			if err != nil {
				return err
			}
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			end, err := c.GetEndpoint(ctx, name)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				utils.JSONPrint(end)
			} else {
				utils.PrintEndpointsTable([]meroxa.Endpoint{*end})
			}
			return nil

		},
	}
}
