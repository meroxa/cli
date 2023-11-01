package secrets

import (
	"context"
	"errors"
	"fmt"

	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
)

type Describe struct {
	args struct {
		nameOrUUID string
	}

	client global.BasicClient
	config config.Config
	logger log.Logger
}

var (
	_ builder.CommandWithBasicClient = (*Describe)(nil)
	_ builder.CommandWithConfig      = (*Describe)(nil)
	_ builder.CommandWithDocs        = (*Describe)(nil)
	_ builder.CommandWithExecute     = (*Describe)(nil)
	_ builder.CommandWithLogger      = (*Describe)(nil)
	_ builder.CommandWithArgs        = (*Describe)(nil)
)

func (*Describe) Usage() string {
	return "describe nameOrUUID"
}

func (*Describe) Docs() builder.Docs {
	return builder.Docs{
		Short: "Describe a Turbine Secret",
		Long: `This command will describe a turbine secret by id or name.
`,
		Example: `meroxa secrets describe nameOrUUID
`,
	}
}

func (d *Describe) Config(cfg config.Config) {
	d.config = cfg
}

func (d *Describe) ParseArgs(args []string) error {
	if len(args) > 0 {
		d.args.nameOrUUID = args[0]
	}
	return nil
}

func (d *Describe) BasicClient(client global.BasicClient) {
	d.client = client
}

func (d *Describe) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Describe) Execute(ctx context.Context) error {
	if d.args.nameOrUUID != "" {
		getSecrets, err := RetrieveSecretsID(ctx, d.client, d.args.nameOrUUID)
		if err != nil {
			return err
		}

		for _, secret := range getSecrets.Items {
			d.logger.Info(ctx, display.PrintTable(secret, displayDetails))
			dashboardURL := fmt.Sprintf("%s/secrets/%s", global.GetMeroxaAPIURL(), secret.ID)
			d.logger.Info(ctx, fmt.Sprintf("\n ✨ To view your secret, visit %s", dashboardURL))
		}
		d.logger.JSON(ctx, getSecrets)
	} else {
		return errors.New("action aborted")
	}
	return nil
}
