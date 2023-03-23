package apps

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine/javascript"
	turbinePY "github.com/meroxa/cli/cmd/meroxa/turbine/python"
	turbineRb "github.com/meroxa/cli/cmd/meroxa/turbine/ruby"
	"github.com/meroxa/cli/log"
)

type Init struct {
	logger     log.Logger
	turbineCLI turbine.CLI
	path       string

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
	return "init APP_NAME [--path pwd] --lang js|go|py"
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
	var (
		err error

		name = i.args.appName
		lang = strings.ToLower(i.flags.Lang)
	)

	name, err = validateAppName(name)
	if err != nil {
		return err
	}

	i.path, err = turbine.GetPath(i.flags.Path)
	if err != nil {
		return err
	}

	i.logger.StartSpinner("\t", fmt.Sprintf("Initializing application %q in %q...", name, i.path))
	if i.turbineCLI == nil {
		i.turbineCLI, err = newTurbineCLI(i.logger, lang, i.path)
		if err != nil {
			i.logger.StopSpinnerWithStatus("\t", log.Failed)
			return err
		}
	}

	if err = i.turbineCLI.Init(ctx, name); err != nil {
		i.logger.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}
	i.logger.StopSpinnerWithStatus("Application directory created!", log.Successful)

	if lang == "go" || lang == turbine.GoLang {
		if err = turbineGo.GoInit(i.logger, i.path+"/"+name, i.flags.SkipModInit, i.flags.ModVendor); err != nil {
			i.logger.StopSpinnerWithStatus("\t", log.Failed)
			return err
		}
	}

	i.logger.StartSpinner("\t", "Running git initialization...")
	if err = i.turbineCLI.GitInit(ctx, i.path+"/"+name); err != nil {
		i.logger.StopSpinnerWithStatus(
			"\tThe final step to 'git init' the Application repo failed. Please complete this step manually.",
			log.Failed)
		return err
	}
	i.logger.StopSpinnerWithStatus("Git initialized successfully!", log.Successful)

	i.logger.Infof(ctx, "Turbine Data Application successfully initialized!\n"+
		"You can start interacting with Meroxa in your app located at \"%s/%s\".\n"+
		"Your Application will not be visible in the Meroxa Dashboard until after deployment.", i.path, name)

	return nil
}

func newTurbineCLI(logger log.Logger, lang, path string) (turbine.CLI, error) {
	switch lang {
	case "go", turbine.GoLang:
		return turbineGo.New(logger, path), nil
	case "js", turbine.JavaScript, turbine.NodeJs:
		return turbineJS.New(logger, path), nil
	case "py", turbine.Python3, turbine.Python:
		return turbinePY.New(logger, path), nil
	case "rb", turbine.Ruby:
		return turbineRb.New(logger, path), nil
	default:
		return nil, fmt.Errorf("language %q not supported. %s", lang, LanguageNotSupportedError)
	}
}
