package apps

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

var (
	_ builder.CommandWithDocs    = (*Init)(nil)
	_ builder.CommandWithFlags   = (*Init)(nil)
	_ builder.CommandWithExecute = (*Init)(nil)
	_ builder.CommandWithLogger  = (*Init)(nil)
)

type Deploy struct {
	logger log.Logger

	flags struct {
		Lang string `long:"lang" short:"l" usage:"language to use (js|go)" required:"true"`
		Path string `long:"path" usage:"path where application will be initialized (current directory as default)"`
	}
}

func (*Deploy) Usage() string {
	return "deploy"
}

func (d *Deploy) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (*Deploy) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Initialize a Meroxa Data Application",
		Example: "meroxa apps deploy my-app",
	}
}

func (d *Deploy) Execute(ctx context.Context) error {
	return nil
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}
