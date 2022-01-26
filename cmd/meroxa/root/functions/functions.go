package functions

import (
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/spf13/cobra"
)

type Functions struct{}

var (
	_ builder.CommandWithAliases     = (*Functions)(nil)
	_ builder.CommandWithDocs        = (*Functions)(nil)
	_ builder.CommandWithFeatureFlag = (*Functions)(nil)
	_ builder.CommandWithSubCommands = (*Functions)(nil)
	_ builder.CommandWithHidden      = (*Functions)(nil)
)

func (*Functions) Usage() string {
	return "functions"
}

func (*Functions) Hidden() bool {
	return true
}

func (*Functions) FeatureFlag() (string, error) {
	return "functions", fmt.Errorf(`no access to the Meroxa functions feature`)
}

func (*Functions) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage functions on Meroxa",
	}
}

func (*Functions) Aliases() []string {
	return []string{"function"}
}

func (*Functions) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Create{}),
		builder.BuildCobraCommand(&List{}),
		builder.BuildCobraCommand(&Describe{}),
	}
}
