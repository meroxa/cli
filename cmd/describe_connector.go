package cmd

import (
	"context"
	"errors"
	"github.com/meroxa/cli/display"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

func DescribeConnectorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connector [name]",
		Short: "Describe connector",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires connector name\n\nUsage:\n  meroxa describe connector <name> [flags]")
			}
			var (
				err  error
				conn *meroxa.Connector
			)
			name := args[0]
			c, err := client()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			conn, err = c.GetConnectorByName(ctx, name)
			if err != nil {
				return err
			}

			if flagRootOutputJSON {
				display.JSONPrint(conn)
			} else {
				display.PrintConnectorsTable([]*meroxa.Connector{conn})
			}
			return nil
		},
	}
}
