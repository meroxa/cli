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

// TODO - implement app init
func (i *Init) Execute(ctx context.Context) error {

	return nil
}
