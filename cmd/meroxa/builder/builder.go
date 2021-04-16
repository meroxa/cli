package builder

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Command interface {
	// Usage is the one-line usage message.
	// Recommended syntax is as follow:
	//   [ ] identifies an optional argument. Arguments that are not enclosed in brackets are required.
	//   ... indicates that you can specify multiple values for the previous argument.
	//   |   indicates mutually exclusive information. You can use the argument to the left of the separator or the
	//       argument to the right of the separator. You cannot use both arguments in a single use of the command.
	//   { } delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are
	//       optional, they are enclosed in brackets ([ ]).
	// Example: add [-F file | -D dir]... [-f format] profile
	Usage() string
}

type CommandWithExecute interface {
	Command
	// Execute is the actual work function. Most commands will implement this.
	Execute(ctx context.Context) error
}

type CommandWithDocs interface {
	Command
	// Docs returns the documentation for the command.
	Docs() Docs
}

// Docs will be shown to the user when typing 'help' as well as in generated docs.
type Docs struct {
	// Short is the short description shown in the 'help' output.
	Short string
	// Long is the long message shown in the 'help <this-command>' output.
	Long string
	// Example is examples of how to use the command.
	Example string
}

type CommandWithAliases interface {
	Command
	// Aliases is an array of aliases that can be used instead of the first word in Usage.
	Aliases() []string
}

type CommandWithArgs interface {
	Command
	// ParseArgs is meant to parse arguments after the command name.
	ParseArgs([]string) error
}

type CommandWithFlags interface {
	Command
	// Flags returns the set of flags on this command.
	Flags() []Flag
}

// Flag describes a single command line flag.
type Flag struct {
	// Long name of the flag.
	Long string
	// Short name of the flag (one character).
	Short string
	// Usage is the description shown in the 'help' output.
	Usage string
	// Required is used to mark the flag as required.
	Required bool
	// Persistent is used to propagate the flag to subcommands.
	Persistent bool
	// Default is the default value when the flag is not explicitly supplied. It should have the same type as the value
	// behind the pointer in field Ptr.
	Default interface{}
	// Ptr is a pointer to the value into which the flag will be parsed.
	Ptr interface{}
}

type CommandWithConfirm interface {
	Command
	// Confirm adds a prompt before the command is executed where the user is asked to write the exact value as
	// wantInput. If the user input matches the command will be executed, otherwise processing will be stopped.
	Confirm(ctx context.Context) (wantInput string)
}

type CommandWithClient interface {
	Command
	// Client provides the meroxa client to the command.
	Client(*meroxa.Client)
}

type CommandWithLogger interface {
	Command
	// Logger provides the logger to the command.
	Logger(log.Logger)
}

type CommandWithSubCommands interface {
	Command
	// SubCommands defines subcommands of a command.
	SubCommands() []*cobra.Command
}

// BuildCobraCommand takes a Command and builds a *cobra.Command from it. It figures out if the command implements any
// other CommandWith* interfaces and configures the cobra command accordingly.
func BuildCobraCommand(c Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: c.Usage(),
	}

	buildCommandWithDocs(cmd, c)
	buildCommandWithAliases(cmd, c)
	buildCommandWithFlags(cmd, c)
	buildCommandWithArgs(cmd, c)
	buildCommandWithLogger(cmd, c)
	buildCommandWithClient(cmd, c)
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
			if f.Default == nil {
				f.Default = ""
			}
			flags.StringVarP(val, f.Long, f.Short, f.Default.(string), f.Usage)
		case *int:
			if f.Default == nil {
				f.Default = 0
			}
			flags.IntVarP(val, f.Long, f.Short, f.Default.(int), f.Usage)
		case *int8:
			if f.Default == nil {
				f.Default = int8(0)
			}
			flags.Int8VarP(val, f.Long, f.Short, f.Default.(int8), f.Usage)
		case *int16:
			if f.Default == nil {
				f.Default = int16(0)
			}
			flags.Int16VarP(val, f.Long, f.Short, f.Default.(int16), f.Usage)
		case *int32:
			if f.Default == nil {
				f.Default = int32(0)
			}
			flags.Int32VarP(val, f.Long, f.Short, f.Default.(int32), f.Usage)
		case *int64:
			if f.Default == nil {
				f.Default = int64(0)
			}
			flags.Int64VarP(val, f.Long, f.Short, f.Default.(int64), f.Usage)
		case *float32:
			if f.Default == nil {
				f.Default = float32(0)
			}
			flags.Float32VarP(val, f.Long, f.Short, f.Default.(float32), f.Usage)
		case *float64:
			if f.Default == nil {
				f.Default = float64(0)
			}
			flags.Float64VarP(val, f.Long, f.Short, f.Default.(float64), f.Usage)
		case *bool:
			if f.Default == nil {
				f.Default = false
			}
			flags.BoolVarP(val, f.Long, f.Short, f.Default.(bool), f.Usage)
		case *time.Duration:
			if f.Default == nil {
				f.Default = time.Duration(0)
			}
			flags.DurationVarP(val, f.Long, f.Short, f.Default.(time.Duration), f.Usage)
		case *[]bool:
			flags.BoolSliceVarP(val, f.Long, f.Short, f.Default.([]bool), f.Usage)
		case *[]float32:
			flags.Float32SliceVarP(val, f.Long, f.Short, f.Default.([]float32), f.Usage)
		case *[]float64:
			flags.Float64SliceVarP(val, f.Long, f.Short, f.Default.([]float64), f.Usage)
		case *[]int32:
			flags.Int32SliceVarP(val, f.Long, f.Short, f.Default.([]int32), f.Usage)
		case *[]int64:
			flags.Int64SliceVarP(val, f.Long, f.Short, f.Default.([]int64), f.Usage)
		case *[]int:
			flags.IntSliceVarP(val, f.Long, f.Short, f.Default.([]int), f.Usage)
		case *[]string:
			flags.StringSliceVarP(val, f.Long, f.Short, f.Default.([]string), f.Usage)
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

func buildCommandWithLogger(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithLogger)
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

		v.Logger(global.NewLogger())
		return nil
	}
}

func buildCommandWithClient(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithClient)
	if !ok {
		return
	}

	old := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
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

	old := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
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

	old := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if old != nil {
			err := old(cmd, args)
			if err != nil {
				return err
			}
		}
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
