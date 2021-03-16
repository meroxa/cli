package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
func UpdateCmd() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a component",
		Long:  `Update a component of the Meroxa platform, including connectors`,
	}

	updateCmd.AddCommand(UpdateConnectorCmd())
	updateCmd.AddCommand(UpdatePipelineCmd())
	updateCmd.AddCommand(UpdateResourceCmd())

	return updateCmd
}

