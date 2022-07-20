package apps

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	turbineCLI "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	turbinejs "github.com/meroxa/cli/cmd/meroxa/turbine_cli/javascript"
	turbinepy "github.com/meroxa/cli/cmd/meroxa/turbine_cli/python"
	"github.com/meroxa/cli/log"
	turbine "github.com/meroxa/turbine-go/init"
)

type Init struct {
	logger log.Logger
	path   string

	args struct {
		appName string
	}

	flags struct {
		Lang        string `long:"lang" short:"l" usage:"language to use (js|go|py)" required:"true"`
		Path        string `long:"path" usage:"path where application will be initialized (current directory as default)"`
		ModVendor   bool   `long:"mod-vendor" usage:"whether to download modules to vendor or globally while initializing a Go application"`
		SkipModInit bool   `long:"skip-mod-init" usage:"whether to run 'go mod init' while initializing a Go application"`
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
	return "init [APP_NAME] [--path pwd] --lang js|go|py"
}

func (*Init) Docs() builder.Docs {
	return builder.Docs{
		Short: "Initialize a Turbine Data Application",
		Example: `meroxa apps init my-app --path ~/code --lang js
meroxa apps init my-app --lang py
meroxa apps init my-app --lang go 			# will be initialized in a dir called my-app in the current directory
meroxa apps init my-app --lang go --skip-mod-init 	# will not initialize the new go module
meroxa apps init my-app --lang go --mod-vendor 		# will initialize the new go module and download dependencies to the vendor directory
meroxa apps init my-app --lang go --path $GOPATH/src/github.com/my.org
`,
		Beta: true,
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

	cmd := exec.Command("git", "config", "--global", "init.defaultBranch", "main")
	cmd.Path = path
	_ = cmd.Run()

	cmd = exec.Command("git", "init", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		i.logger.Error(ctx, string(output))
		return err
	}
	return nil
}

func validateAppName(name string) (string, error) {
	var err error

	// must be lowercase because of reusing the name for the Docker image
	name = strings.ToLower(name)

	// Platform API requires the first character be a letter and
	// that the whole name be alphanumeric with dashes and underscores.
	r := regexp.MustCompile(`^([a-z][a-z0-9-_]*)$`)
	matches := r.FindStringSubmatch(name)
	if len(matches) == 0 {
		err = fmt.Errorf(
			"invalid application name: %s;"+
				" should start with a letter, be alphanumeric, and only have dashes as separators",
			name)
	}
	return name, err
}

func (i *Init) Execute(ctx context.Context) error {
	name := i.args.appName
	lang := i.flags.Lang

	name, err := validateAppName(name)
	if err != nil {
		return err
	}

	i.path, err = turbineCLI.GetPath(i.flags.Path)
	if err != nil {
		return err
	}

	i.logger.StartSpinner("\t", fmt.Sprintf("Initializing application %q in %q...", name, i.path))
	switch lang {
	case "go", GoLang:
		err = turbine.Init(name, i.path)
		if err != nil {
			i.logger.StopSpinnerWithStatus("\t", log.Failed)
			return err
		}
		i.logger.StopSpinnerWithStatus("Application directory created!", log.Successful)
		err = turbineCLI.GoInit(i.logger, i.path+"/"+name, i.flags.SkipModInit, i.flags.ModVendor)
	case "js", JavaScript, NodeJs:
		err = turbinejs.Init(ctx, i.logger, name, i.path)
	case "py", Python3, Python:
		err = turbinepy.Init(ctx, i.logger, name, i.path)
	default:
		i.logger.StopSpinnerWithStatus("\t", log.Failed)
		return fmt.Errorf("language %q not supported. %s", lang, LanguageNotSupportedError)
	}
	if err != nil {
		i.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	if lang != "go" && lang != GoLang {
		i.logger.StopSpinnerWithStatus("Application directory created!", log.Successful)
	}
	i.logger.StartSpinner("\t", "Running git initialization...")
	err = i.GitInit(ctx, i.path+"/"+name)
	if err != nil {
		i.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	i.logger.StopSpinnerWithStatus("Git initialized successfully!", log.Successful)
	i.logger.Infof(ctx, "Turbine Data Application successfully initialized!\n"+
		"You can start interacting with Meroxa in your app located at \"%s/%s\".\n"+
		"Your Application will not be visible in the Meroxa Dashboard until after deployment.", i.path, name)

	return nil
}
