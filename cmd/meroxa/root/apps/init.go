package apps

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Init struct {
	logger log.Logger

	args struct {
		appName string
	}

	flags struct {
		Lang string `long:"lang" short:"l" usage:"language to use (js|go)" required:"true"`
		Path string `long:"path" usage:"path where application will be initialized (current directory as default)"`
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
	return "init [APP_NAME] [--path pwd] --lang js|go"
}

func (*Init) Docs() builder.Docs {
	return builder.Docs{
		Short: "Initialize a Meroxa Data Application",
		Example: "meroxa apps init my-app --path ~/code --lang js" +
			"meroxa apps init my-app --lang go # will be initialized in current directory",
	}
}

func (i *Init) Flags() []builder.Flag {
	return builder.BuildFlags(&i.flags)
}

func (i *Init) ParseArgs(args []string) error {
	// TODO: generate app's name in behalf of the user
	if len(args) < 1 {
		return errors.New("requires an application name")
	}

	i.args.appName = args[0]
	return nil
}

func (i *Init) Execute(ctx context.Context) error {
	name := i.args.appName
	lang := i.flags.Lang

	switch lang {
	case "go", "golang":
		// TODO: Implement apps init for go.
	case "js", "javascript", "nodejs":
		cmd := exec.Command("npx", "turbine", "generate", name)
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		i.logger.Info(ctx, string(stdout))
	default:
		return fmt.Errorf("language %q not supported. Currently, we support \"javascript\" and \"go\"", lang)
	}

	return nil
}

func (i *Init) Logger(logger log.Logger) {
	i.logger = logger
}
