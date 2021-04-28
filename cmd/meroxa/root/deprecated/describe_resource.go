package deprecated

import (
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

// DescribeResourceCmd represents the `meroxa describe resource` command.
func DescribeResourceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resource NAME",
		Short: "Describe resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires resource name\n\nUsage:\n  meroxa describe resource NAME [flags]")
			}
			name := args[0]

			c, err := global.NewClient()
			if err != nil {
				return err
			}

			res, err := c.GetResourceByName(cmd.Context(), name)
			if err != nil {
				return err
			}

			if FlagRootOutputJSON {
				utils.JSONPrint(res)
			} else {
				utils.PrintResourcesTable([]*meroxa.Resource{res})
			}
			return nil
		},
	}
}