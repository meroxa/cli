package secrets

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
)

type Create struct {
	flags struct {
		Data string `long:"data" usage:"Secret's data, passed as a JSON string"`
	}
	args struct {
		secretName string
	}
	client global.BasicClient
	config config.Config
	logger log.Logger
}

var (
	_ builder.CommandWithBasicClient = (*Create)(nil)
	_ builder.CommandWithConfig      = (*Create)(nil)
	_ builder.CommandWithDocs        = (*Create)(nil)
	_ builder.CommandWithExecute     = (*Create)(nil)
	_ builder.CommandWithFlags       = (*Create)(nil)
	_ builder.CommandWithLogger      = (*Create)(nil)
	_ builder.CommandWithArgs        = (*Create)(nil)
)

func (*Create) Usage() string {
	return "create NAME --data '{}'"
}

func (*Create) Docs() builder.Docs {
	return builder.Docs{
		Short: "Create a Conduit Secret",
		Long: `This command will create a secret as promted by the user.'
After successful creation, the secret can be used in a connector. 
`,
		Example: `meroxa secret create NAME
		          meroxa secret create NAME --data '{}'
		`,
	}
}

func (d *Create) ParseArgs(args []string) error {
	if len(args) > 0 {
		d.args.secretName = args[0]
	}
	return nil
}

func (d *Create) Config(cfg config.Config) {
	d.config = cfg
}

func (d *Create) BasicClient(client global.BasicClient) {
	d.client = client
}

func (d *Create) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Create) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Create) Execute(ctx context.Context) error {
	if d.flags.Data == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("To proceed, enter the secret's data as a JSON string: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if len(input) == 0 {
			return errors.New("action aborted")
		}
		d.flags.Data = strings.TrimRight(input, "\r\n")
	}

	secret := &Secrets{
		Name: d.args.secretName,
	}
	err := json.Unmarshal([]byte(d.flags.Data), &secret.Data)
	if err != nil {
		return err
	}

	d.logger.Infof(ctx, "Adding a secret %q...", d.args.secretName)
	response, err := d.client.CollectionRequest(ctx, "POST", collectionName, "", secret, nil)
	if err != nil {
		return err
	}

	responseSecret := Secrets{}
	err = json.NewDecoder(response.Body).Decode(&responseSecret)
	if err != nil {
		return err
	}

	d.logger.Infof(ctx, "Secret %q successfully added", responseSecret.Name)
	d.logger.JSON(ctx, responseSecret)
	dashboardURL := fmt.Sprintf("\n ✨ To view your secrets, visit %s/secrets/%s", global.GetMeroxaAPIURL(), responseSecret.ID)
	d.logger.Info(ctx, dashboardURL)

	return nil
}
