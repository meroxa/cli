package apps

import (
	"context"
	"errors"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/meroxa/cli/log"
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
	return "init APP_NAME [--path pwd] --lang js|go|py"
}

func (*Init) Docs() builder.Docs {
	return builder.Docs{
		Short: "Initialize a Conduit Data Application",
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

func (i *Init) Execute(_ context.Context) error {
	var err error

	// name, err = validateAppName(name)
	// if err != nil {
	// 	return err
	// }

	i.path, err = GetPath(i.flags.Path)
	if err != nil {
		return err
	}

	// i.logger.StartSpinner("\t", fmt.Sprintf("Initializing application %q in %q...", name, i.path))
	// if i.turbineCLI == nil {
	// 	i.turbineCLI, err = newTurbineCLI(i.logger, lang, i.path)
	// 	if err != nil {
	// 		i.logger.StopSpinnerWithStatus("\t", log.Failed)
	// 		return err
	// 	}
	// }

	// if err = i.turbineCLI.Init(ctx, name); err != nil {
	// 	i.logger.StopSpinnerWithStatus("\t", log.Failed)
	// 	return err
	// }
	// i.logger.StopSpinnerWithStatus("Application directory created!", log.Successful)

	// if lang == "go" || lang == string(ir.GoLang) {
	// 	if err = turbineGo.GoInit(i.logger, i.path+"/"+name, i.flags.SkipModInit, i.flags.ModVendor); err != nil {
	// 		i.logger.StopSpinnerWithStatus("\t", log.Failed)
	// 		return err
	// 	}
	// }

	// i.logger.StartSpinner("\t", "Running git initialization...")
	// if err = i.turbineCLI.GitInit(ctx, i.path+"/"+name); err != nil {
	// 	i.logger.StopSpinnerWithStatus(
	// 		"\tThe final step to 'git init' the Application repo failed. Please complete this step manually.",
	// 		log.Failed)
	// 	return err
	// }
	// i.logger.StopSpinnerWithStatus("Git initialized successfully!", log.Successful)

	// appPath := filepath.Join(i.path, name)

	// i.logger.Infof(ctx, "Conduit Data Application successfully initialized!\n"+
	// 	"You can start interacting with Meroxa in your app located at \"%s\".\n"+
	// 	"Your Application will not be visible in the Meroxa Dashboard until after deployment.", appPath)

	return nil
}
