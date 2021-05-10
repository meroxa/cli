package deprecated

import (
	"github.com/spf13/cobra"
)

func RegisterCommands(cmd *cobra.Command) {
	cmd.AddCommand(addCmd())
	cmd.AddCommand(createCmd())
	cmd.AddCommand(describeCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(logsCmd())

	// To migrate
	cmd.AddCommand((&Remove{}).Command())
	cmd.AddCommand(UpdateCmd())
}
