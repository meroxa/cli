package secrets

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
)

type Remove struct {
	flags struct {
		Force bool `long:"force" short:"f" default:"false" usage:"skip confirmation"`
	}
	args struct {
		nameOrUUID string
	}

	client global.BasicClient
	config config.Config
	logger log.Logger
}

var (
	_ builder.CommandWithBasicClient = (*Remove)(nil)
	_ builder.CommandWithConfig      = (*Remove)(nil)
	_ builder.CommandWithDocs        = (*Remove)(nil)
	_ builder.CommandWithExecute     = (*Remove)(nil)
	_ builder.CommandWithFlags       = (*Remove)(nil)
	_ builder.CommandWithLogger      = (*Remove)(nil)
	_ builder.CommandWithArgs        = (*Remove)(nil)
	_ builder.CommandWithAliases     = (*Remove)(nil)
)

func (*Remove) Usage() string {
	return "remove [--path pwd]"
}

func (*Remove) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Remove a Turbine Secret",
		Long:    `This command will remove the secret specified either by name or id`,
		Example: `meroxa apps remove nameOrUUID`,
	}
}

func (d *Remove) Config(cfg config.Config) {
	d.config = cfg
}

func (d *Remove) Aliases() []string {
	return []string{"rm", "delete"}
}

func (d *Remove) BasicClient(client global.BasicClient) {
	d.client = client
}

func (d *Remove) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Remove) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Remove) ParseArgs(args []string) error {
	if len(args) > 0 {
		d.args.nameOrUUID = args[0]
	}
	return nil
}

func (d *Remove) Execute(ctx context.Context) error {
	if !d.flags.Force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("To proceed, type %q or re-run this command with --force\n▸ ", d.args.nameOrUUID)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if d.args.nameOrUUID != strings.TrimRight(input, "\r\n") {
			return errors.New("action aborted")
		}
	}

	getSecrets, err := RetrieveSecretsID(ctx, d.client, d.args.nameOrUUID)
	if err != nil {
		return err
	}

	d.logger.Infof(ctx, "Removing secret %q...", d.args.nameOrUUID)
	_, err = d.client.CollectionRequest(ctx, "DELETE", collectionName, getSecrets.Items[0].ID, nil, nil)
	if err != nil {
		return err
	}

	d.logger.Infof(ctx, "Secret %q successfully removed", d.args.nameOrUUID)

	return nil
}
