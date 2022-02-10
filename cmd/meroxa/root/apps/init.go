package apps

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Init struct {
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Lang string `long:"lang" short:"l" usage:"language to use (go | javascript)" required:"true"`
	}
}

var (
	_ builder.CommandWithDocs    = (*Init)(nil)
	_ builder.CommandWithArgs    = (*Init)(nil)
	_ builder.CommandWithFlags   = (*Init)(nil)
	_ builder.CommandWithExecute = (*Init)(nil)
	_ builder.CommandWithLogger  = (*Init)(nil)
)

func (*Init) Usage() string {
	return "init"
}

func (*Init) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Initialize a Meroxa Data Application",
		Example: "meroxa apps init",
	}
}

func (i *Init) Flags() []builder.Flag {
	return builder.BuildFlags(&i.flags)
}

func (i *Init) ParseArgs(args []string) error {
	if len(args) > 0 {
		i.args.Name = args[0]
	}
	return nil
}

func (i *Init) Execute(ctx context.Context) error {
	name := i.args.Name
	lang := i.flags.Lang

	if lang == "javascript" {
		cmd := exec.Command("npx", "turbine", "generate", name)
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		i.logger.Info(ctx, string(stdout))
	} else if lang == "go" {
		// TODO: Implement apps init for go.
	} else {
		return fmt.Errorf("unsupported language: %s", lang)
	}

	return nil
}

func (i *Init) Logger(logger log.Logger) {
	i.logger = logger
}
