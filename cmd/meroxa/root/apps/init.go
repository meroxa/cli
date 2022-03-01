package apps

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
	turbine "github.com/meroxa/turbine/init"
)

type Init struct {
	logger log.Logger
	path   string

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

func (i *Init) Logger(logger log.Logger) {
	i.logger = logger
}

func (i *Init) Flags() []builder.Flag {
	return builder.BuildFlags(&i.flags)
}

func (i *Init) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires an application name")
	}

	i.args.appName = args[0]
	return nil
}

func (i *Init) GitInit(ctx context.Context, path string) error {
	if path == "" {
		return errors.New("path is required")
	}

	cmd := exec.Command("git", "init", path)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (i *Init) Execute(ctx context.Context) error {
	name := i.args.appName
	lang := i.flags.Lang

	i.path = turbineCLI.GetPath(i.flags.Path)

	i.logger.Infof(ctx, "Initializing application %q in %q...", name, i.path)
	switch lang {
	case "go", GoLang:
		err := turbine.Init(i.path, name)
		if err != nil {
			return err
		}
		i.logger.Infof(ctx, "Application successfully initialized!\n"+
			"You can start interacting with Meroxa in your app located at \"%s/%s\"", i.path, name)
	case "js", JavaScript, NodeJs:
		cmd := exec.Command("npx", "turbine", "generate", i.path, name)
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		i.logger.Info(ctx, string(stdout))
	default:
		return fmt.Errorf("language %q not supported. Currently, we support \"javascript\" and \"go\"", lang)
	}

	err := i.GitInit(ctx, i.path+"/"+name)
	if err != nil {
		return err
	}

	return nil
}