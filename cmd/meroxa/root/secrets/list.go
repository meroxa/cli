package secrets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
)

type List struct {
	client global.BasicClient
	config config.Config
	logger log.Logger
}

var (
	_ builder.CommandWithBasicClient = (*List)(nil)
	_ builder.CommandWithConfig      = (*List)(nil)
	_ builder.CommandWithDocs        = (*List)(nil)
	_ builder.CommandWithExecute     = (*List)(nil)
	_ builder.CommandWithLogger      = (*List)(nil)
	_ builder.CommandWithAliases     = (*List)(nil)
)

func (*List) Usage() string {
	return "list [--path pwd]"
}

func (*List) Docs() builder.Docs {
	return builder.Docs{
		Short: "List all Conduit Secrets",
		Long: `This command will list all the secrets defined on the platform.
`,
		Example: `meroxa secrets list`,
	}
}

func (d *List) Aliases() []string {
	return []string{"ls"}
}

func (d *List) Config(cfg config.Config) {
	d.config = cfg
}

func (d *List) BasicClient(client global.BasicClient) {
	d.client = client
}

func (d *List) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *List) Execute(ctx context.Context) error {
	secrets := &ListSecrets{}

	response, err := d.client.CollectionRequest(ctx, "GET", collectionName, "", nil, nil)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&secrets)
	if err != nil {
		return err
	}

	d.logger.Info(ctx, display.PrintList(secrets.Items, displayDetails))
	d.logger.JSON(ctx, secrets)
	output := fmt.Sprintf("\n ✨ To view your secrets, visit %s/secrets", global.GetMeroxaAPIURL())
	d.logger.Info(ctx, output)

	return nil
}
