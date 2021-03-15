package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a component",
	Long:  `Update a component of the Meroxa platform, including connectors`,
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
