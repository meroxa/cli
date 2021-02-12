package cmd

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(logsCmd)
	logsCmd.AddCommand(logsConnectorCmd)
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print logs for a component",
}

var logsConnectorCmd = &cobra.Command{
	Use:   "connector <name>",
	Short: "Print logs for a connector",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		_, err = io.Copy(os.Stderr, resp.Body)

		return err
	},
}
