/*
Copyright © 2021 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package builder

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/cased/cased-go"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

type CommandWithClient interface {
	Command
	// Client provides the meroxa client to the command.
	Client(*meroxa.Client)
}

type CommandWithConfig interface {
	Command
	Config(config.Config)
}

type CommandWithConfirm interface {
	Command
	// Confirm adds a prompt before the command is executed where the user is asked to write the exact value as
	// wantInput. If the user input matches the command will be executed, otherwise processing will be stopped.
	Confirm(ctx context.Context) (wantInput string)
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

type CommandWithExecute interface {
	Command
	// Execute is the actual work function. Most commands will implement this.
	Execute(ctx context.Context) error
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
	// Hidden is used to mark the flag as hidden.
	Hidden bool
}

type CommandWithHidden interface {
	Command
	// Hidden returns the desired hidden value for the command.
	Hidden() bool
}

type CommandWithLogger interface {
	Command
	// Logger provides the logger to the command.
	Logger(log.Logger)
}

type CommandWithEvent interface {
	Command
	Event() cased.AuditEvent
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

	buildCommandWithAliases(cmd, c)
	buildCommandWithArgs(cmd, c)
	buildCommandWithClient(cmd, c)
	buildCommandWithConfig(cmd, c)
	buildCommandWithConfirm(cmd, c)
	buildCommandWithDocs(cmd, c)
	buildCommandWithExecute(cmd, c)
	buildCommandWithFlags(cmd, c)
	buildCommandWithHidden(cmd, c)
	buildCommandWithLogger(cmd, c)
	buildCommandWithSubCommands(cmd, c)

	// this has to be the last function so it captures all errors from RunE
	buildCommandWithEvent(cmd, c)

	return cmd
}

func buildCommandWithAliases(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithAliases)
	if !ok {
		return
	}

	cmd.Aliases = v.Aliases()
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

func buildCommandWithConfig(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithConfig)
	if !ok {
		return
	}

	// Inject global.Config.
	oldPreRunE := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if oldPreRunE != nil {
			err := oldPreRunE(cmd, args)
			if err != nil {
				return err
			}
		}

		v.Config(global.Config)
		return nil
	}

	// Make sure writes on file in the end.
	oldPostRunE := cmd.PostRunE
	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		if oldPostRunE != nil {
			err := oldPostRunE(cmd, args)
			if err != nil {
				return err
			}
		}

		err := global.Config.WriteConfig()

		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				err = global.Config.SafeWriteConfig()
			}
			if err != nil {
				return fmt.Errorf("meroxa: could not write config file: %v", err)
			}
		}

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

		// do not prompt for confirmation when --force (or --yolo 😜) is set
		if force || yolo {
			return nil
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

func getCLIUserInfo() (actor, actorUUID string, err error) {
	// Require login
	_, _, err = global.RequireLogin()

	/*
		 	We don't report client issues to the customer as it'll likely require `meroxa login` for any command.
			There are command that don't require client such as `meroxa env`, and we wouldn't like to throw an error,
			just because we can't emit events.
	*/
	if err != nil {
		return "", "", nil
	}

	// fetch actor account.
	actor = global.Config.GetString("MEROXA_ACTOR")
	actorUUID = global.Config.GetString("MEROXA_ACTOR_UUID")

	if actor == "" || actorUUID == "" {
		// call api to fetch
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // nolint:gomnd
		defer cancel()

		m, err := global.NewClient()

		if err != nil {
			return "", "", fmt.Errorf("meroxa: could not create Meroxa client: %v", err)
		}

		account, err := m.GetUser(ctx)

		if err != nil {
			return "", "", fmt.Errorf("meroxa: could not fetch Meroxa user: %v", err)
		}

		actor = account.Email
		actorUUID = account.UUID

		global.Config.Set("MEROXA_ACTOR", actor)
		global.Config.Set("MEROXA_ACTOR_UUID", actorUUID)

		err = global.Config.WriteConfig()

		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				err = global.Config.SafeWriteConfig()
			}
			if err != nil {
				return "", "", fmt.Errorf("meroxa: could not write config file: %v", err)
			}
		}
	}

	return actor, actorUUID, nil
}

// This runs for all commands.
func buildCommandWithEvent(cmd *cobra.Command, c Command) {
	oldRunE := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if oldRunE != nil {
			err := oldRunE(cmd, args)
			if err != nil {
				return err
			}

			// This will be empty on logout, maybe do it before?
			actor, actorUUID, err := getCLIUserInfo()

			if err != nil {
				return err
			}

			// build our event
			event := cased.AuditEvent{
				"timestamp":  time.Now().UTC(),
				"user_agent": fmt.Sprintf("meroxa/%s %s/%s", global.Version, runtime.GOOS, runtime.GOARCH),
			}

			if actor != "" {
				event["actor"] = actor
			}

			if actorUUID != "" {
				event["actor_uuid"] = actorUUID
			}

			var action string

			// TODO: Implement something that could look up all the way up until meroxa (meroxa create resources...)
			// something like it determines how many levels since root and then until current cmd
			if cmd.HasParent() {
				if cmd.Parent().HasParent() {
					action = fmt.Sprintf("%s.%s.%s", cmd.Parent().Parent().Use, cmd.Parent().Use, cmd.Use)
				} else {
					action = fmt.Sprintf("%s.%s", cmd.Parent().Use, cmd.Use)
				}
			} else {
				action = cmd.Use
			}

			if cmd.Use != cmd.CalledAs() {
				event["command.alias"] = cmd.CalledAs()
			}

			if len(args) > 0 {
				event["command.args"] = args
			}

			event["action"] = action
			event["use"] = action

			if err != nil {
				event["error"] = err
			}

			if cmd.HasFlags() {
				cmd.Flags().Visit(func(flag *pflag.Flag) {
					event["command.flags"] = flag.Name
				})
			}

			if cmd.Deprecated != "" {
				event["command.deprecated"] = "true"
			}

			v, ok := c.(CommandWithEvent)
			if ok {
				metadata := v.Event()

				// merge default event with what's defined in the command.
				for k, v := range metadata {
					event[k] = v
				}
			}

			//fmt.Println(event)

			casedAPIKey := global.Config.GetString("CASED_API_KEY")
			publisher := global.NewPublisher(casedAPIKey)
			cased.SetPublisher(publisher)
			err = cased.Publish(event)
			if err != nil {
				return fmt.Errorf("meroxa: couldn't emit audit trail event: %v", err)
			}
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

// nolint:funlen,gocyclo // this function has a big switch statement, can't get around that
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
			if f.Default == nil {
				f.Default = []bool(nil)
			}
			flags.BoolSliceVarP(val, f.Long, f.Short, f.Default.([]bool), f.Usage)
		case *[]float32:
			if f.Default == nil {
				f.Default = []float32(nil)
			}
			flags.Float32SliceVarP(val, f.Long, f.Short, f.Default.([]float32), f.Usage)
		case *[]float64:
			if f.Default == nil {
				f.Default = []float64(nil)
			}
			flags.Float64SliceVarP(val, f.Long, f.Short, f.Default.([]float64), f.Usage)
		case *[]int32:
			if f.Default == nil {
				f.Default = []int32(nil)
			}
			flags.Int32SliceVarP(val, f.Long, f.Short, f.Default.([]int32), f.Usage)
		case *[]int64:
			if f.Default == nil {
				f.Default = []int64(nil)
			}
			flags.Int64SliceVarP(val, f.Long, f.Short, f.Default.([]int64), f.Usage)
		case *[]int:
			if f.Default == nil {
				f.Default = []int(nil)
			}
			flags.IntSliceVarP(val, f.Long, f.Short, f.Default.([]int), f.Usage)
		case *[]string:
			if f.Default == nil {
				f.Default = []string(nil)
			}
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

		if f.Hidden {
			err := flags.MarkHidden(f.Long)
			if err != nil {
				panic(fmt.Errorf("could not mark flag hidden: %w", err))
			}
		}
	}
}

func buildCommandWithHidden(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithHidden)
	if !ok {
		return
	}

	cmd.Hidden = v.Hidden()
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

func buildCommandWithSubCommands(cmd *cobra.Command, c Command) {
	v, ok := c.(CommandWithSubCommands)
	if !ok {
		return
	}

	for _, sub := range v.SubCommands() {
		cmd.AddCommand(sub)
	}
}
