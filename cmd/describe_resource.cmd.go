package cmd

import (
	"context"
	"errors"
	"github.com/meroxa/cli/display"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

func DescribeResourceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resource <name>",
		Short: "Describe resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires resource name\n\nUsage:\n  meroxa describe resource <name> [flags]")
			}
			name := args[0]

			c, err := client()
			if err != nil {
				return err
			}
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			res, err := c.GetResourceByName(ctx, name)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				display.JSONPrint(res)
			} else {
				display.PrintResourcesTable([]*meroxa.Resource{res})
			}
			return nil
		},
	}
}
