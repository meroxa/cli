package fixtures

import (
	"fmt"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/spf13/cobra"
)

type Fixtures struct{}

var (
	_ builder.CommandWithAliases     = (*Fixtures)(nil)
	_ builder.CommandWithDocs        = (*Fixtures)(nil)
	_ builder.CommandWithFeatureFlag = (*Fixtures)(nil)
	_ builder.CommandWithSubCommands = (*Fixtures)(nil)
	_ builder.CommandWithHidden      = (*Fixtures)(nil)
)

func (*Fixtures) Usage() string {
	return "fixtures"
}

func (*Fixtures) Hidden() bool {
	return true
}

func (*Fixtures) FeatureFlag() (string, error) {
	return "fixtures", fmt.Errorf(`no access to the Meroxa fixtures feature`)
}

func (*Fixtures) Docs() builder.Docs {
	return builder.Docs{
		Short: "Managed Turbine fixtures",
	}
}

func (*Fixtures) Aliases() []string {
	return []string{"fixture"}
}

func (*Fixtures) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Fetch{}),
	}
}
