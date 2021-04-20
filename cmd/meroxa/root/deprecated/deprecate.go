package deprecated

import (
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/root/old"
	"github.com/meroxa/cli/cmd/meroxa/root/resource"
	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:        "add",
		Deprecated: "use `[connectors | endpoints | pipelines | resources] create` instead",
	}
}

func addResourceCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&resource.CreateResource{})
	cmd.Use = "resource"
	cmd.Deprecated = "use `resource create` instead"
	return cmd
}

func RegisterCommands(cmd *cobra.Command) {
	// meroxa add resource
	addCmd := addCmd()
	addCmd.AddCommand(addResourceCmd())
	cmd.AddCommand(addCmd)

	cmd.AddCommand(old.CreateCmd())
	cmd.AddCommand(old.DescribeCmd())
	cmd.AddCommand(old.ListCmd())
	cmd.AddCommand(old.LogsCmd())
	cmd.AddCommand(old.OpenCmd())
	cmd.AddCommand((&old.Remove{}).Command())
	cmd.AddCommand(old.UpdateCmd())
}
