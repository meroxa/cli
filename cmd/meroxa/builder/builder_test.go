package builder_test

import (
	"reflect"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/spf13/cobra"
)

type testCmd struct {
	flagLongFoo string
}

var (
	_ builder.CommandWithDocs        = (*testCmd)(nil)
	_ builder.CommandWithAliases     = (*testCmd)(nil)
	_ builder.CommandWithFlags       = (*testCmd)(nil)
	_ builder.CommandWithSubCommands = (*testCmd)(nil)
	// _ builder.CommandWithArgs        = (*testCmd)(nil)
	// _ builder.CommandWithLogger      = (*testCmd)(nil)
	// _ builder.CommandWithClient      = (*testCmd)(nil)
	// _ builder.CommandWithConfirm     = (*testCmd)(nil)
	// _ builder.CommandWithExecute     = (*testCmd)(nil)
)

func (c *testCmd) Usage() string {
	return "cmd1"
}
func (c *testCmd) Docs() builder.Docs {
	return builder.Docs{
		Short:   "short-foo",
		Long:    "long-bar",
		Example: "example-baz",
	}
}
func (c *testCmd) Aliases() []string {
	return []string{"foo", "bar"}
}
func (c *testCmd) Flags() []builder.Flag {
	return []builder.Flag{
		{Long: "long-foo", Short: "l", Usage: "test flag", Required: false, Persistent: false, Ptr: &c.flagLongFoo},
	}
}
func (c *testCmd) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&subCmd{}),
	}
}

type subCmd struct{}

func (c *subCmd) Usage() string {
	return "cmd2"
}

func TestBuildCobraCommand_OnlyUsage(t *testing.T) {
	cmd := &testCmd{}

	want := &cobra.Command{
		Use:     "cmd1",
		Aliases: []string{"foo", "bar"},
		Short:   "short-foo",
		Long:    "long-bar",
		Example: "example-baz",
	}
	want.Flags().StringVarP(&cmd.flagLongFoo, "long-foo", "l", "", "test flag")
	want.AddCommand(&cobra.Command{Use: "cmd2"})

	got := builder.BuildCobraCommand(cmd)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`expected "%v", got "%v"`, want, got)
	}
}
