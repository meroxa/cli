package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/meroxa/cli/display"
	"github.com/spf13/cobra"
)

var (
	state string // connector state
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a component",
	Long:  `Update a component of the Meroxa platform, including connectors`,
}

var updateConnectorCmd = &cobra.Command{
	Use:   "connector <name> --state <pause|resume|restart>",
	Short: "Update connector state",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires connector name\n\nUsage:\n  meroxa update connector <name> --state <state>")
		}

		// Connector Name
		conName := args[0]

		c, err := client()
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// call meroxa-go to update connector status with name
		if !flagRootOutputJSON {
			fmt.Printf("Updating %s connector...\n", conName)
		}

		con, err := c.UpdateConnectorStatus(ctx, conName, state)
		if err != nil {
			return err
		}

		if flagRootOutputJSON {
			display.JSONPrint(con)
		} else {
			fmt.Printf("Connector %s successfully updated!\n", con.Name)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	// Subcommands
	updateCmd.AddCommand(updateConnectorCmd)
	updateConnectorCmd.Flags().StringVarP(&state, "state", "", "", "connector state")
	updateConnectorCmd.MarkFlagRequired("state")
}
