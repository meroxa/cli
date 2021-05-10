package deprecated

import (
	"github.com/spf13/cobra"
)

func RegisterCommands(cmd *cobra.Command) {
	cmd.AddCommand(addCmd()) // ✅.

	cmd.AddCommand(CreateCmd())
	cmd.AddCommand(DescribeCmd())
	cmd.AddCommand(ListCmd())
	cmd.AddCommand(LogsCmd())
	cmd.AddCommand((&Remove{}).Command())
	cmd.AddCommand(UpdateCmd())
}
