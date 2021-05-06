package deprecated

import (
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/root/deprecated/add"
	"github.com/spf13/cobra"
)

func RegisterCommands(cmd *cobra.Command) {
	cmd.AddCommand(builder.BuildCobraCommand(&add.Add{}))
	cmd.AddCommand(CreateCmd())
	cmd.AddCommand(DescribeCmd())
	cmd.AddCommand(ListCmd())
	cmd.AddCommand(LogsCmd())
	cmd.AddCommand((&Remove{}).Command())
	cmd.AddCommand(UpdateCmd())
}
