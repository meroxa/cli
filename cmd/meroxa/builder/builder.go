package builder

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Command interface {
	Usage() string
}

type CommandWithExecute interface {
	Command
	Execute(ctx context.Context) error
}

type CommandWithDocs interface {
	Command
	Docs() Docs
}

type Docs struct {
	Short   string
	Long    string
	Example string
}

type CommandWithAliases interface {
	Command
	Aliases() []string
}

type CommandWithArgs interface {
	Command
	ParseArgs([]string) error
}

type CommandWithFlags interface {
	Command
	Flags() []Flag
}

type Flag struct {
	Long       string
	Short      string
	Usage      string
	Default    string
	Required   bool
	Persistent bool
	Ptr        interface{}
}

type CommandWithConfirm interface {
	Command
	Confirm(ctx context.Context) (wantInput string)
}

type CommandWithClient interface {
	Command
	Client(*meroxa.Client)
}

type CommandWithSubCommands interface {
	Command
	SubCommands() []*cobra.Command
}

func BuildCobraCommand(c Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: c.Usage(),
	}

	buildCommandWithDocs(cmd, c)
	buildCommandWithAliases(cmd, c)
	buildCommandWithClient(cmd, c)
	buildCommandWithFlags(cmd, c)
	buildCommandWithArgs(cmd, c)
	buildCommandWithConfirm(cmd, c)
	buildCommandWithExecute(cmd, c)
	buildCommandWithSubCommands(cmd, c)

	return cmd
}

func buildCommandWithDocs(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithDocs)
	if !ok {
		return
	}

	docs := v.Docs()
	cmd.Long = docs.Long
	cmd.Short = docs.Short
	cmd.Example = docs.Example
}

func buildCommandWithAliases(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithAliases)
	if !ok {
		return
	}

	cmd.Aliases = v.Aliases()
}

func buildCommandWithClient(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithClient)
	if !ok {
		return
	}

	old := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if old != nil {
			err := old(cmd, args)
			if err != nil {
				return err
			}
		}
		c, err := global.NewClient()
		if err != nil {
			return err
		}
		v.Client(c)
		return nil
	}
}

func buildCommandWithFlags(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithFlags)
	if !ok {
		return
	}

	for _, f := range v.Flags() {
		var flags *pflag.FlagSet
		if f.Persistent {
			flags = cmd.PersistentFlags()
		} else {
			flags = cmd.Flags()
		}

		switch val := f.Ptr.(type) {
		case *string:
			flags.StringVarP(val, f.Long, f.Short, f.Default, f.Usage)
		// TODO add more types
		default:
			panic(fmt.Errorf("unexpected flag value type: %T", val))
		}

		if f.Required {
			err := cobra.MarkFlagRequired(flags, f.Long)
			if err != nil {
				panic(fmt.Errorf("could not mark flag required: %w", err))
			}
		}
	}
}

func buildCommandWithArgs(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithArgs)
	if !ok {
		return
	}

	old := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if old != nil {
			err := old(cmd, args)
			if err != nil {
				return err
			}
		}
		return v.ParseArgs(args)
	}
}

func buildCommandWithConfirm(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithConfirm)
	if !ok {
		return
	}

	var (
		force bool
		yolo  bool
	)
	cmd.Flags().BoolVarP(&force, "force", "f", false, "skip confirmation")
	cmd.Flags().BoolVarP(&yolo, "yolo", "", false, "skip confirmation")
	err := cmd.Flags().MarkHidden("yolo")
	if err != nil {
		panic(fmt.Errorf("could not mark flag hidden: %w", err))
	}

	old := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if old != nil {
			err := old(cmd, args)
			if err != nil {
				return err
			}
		}
		wantInput := v.Confirm(cmd.Context())

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("To proceed, type %q or re-run this command with --force\n▸ ", wantInput)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if wantInput != strings.TrimSuffix(input, "\n") {
			return errors.New("action aborted")
		}

		return nil
	}
}

func buildCommandWithExecute(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithExecute)
	if !ok {
		return
	}

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		return v.Execute(cmd.Context())
	}
}

func buildCommandWithSubCommands(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithSubCommands)
	if !ok {
		return
	}

	for _, sub := range v.SubCommands() {
		cmd.AddCommand(sub)
	}
}
