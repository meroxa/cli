package secrets

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
)

type Update struct {
	args struct {
		nameOrUUID string
	}

	flags struct {
		Data string `long:"data" usage:"Secret's data, passed as a JSON string"`
	}

	client global.BasicClient
	config config.Config
	logger log.Logger
}

var (
	_ builder.CommandWithBasicClient = (*Update)(nil)
	_ builder.CommandWithConfig      = (*Update)(nil)
	_ builder.CommandWithDocs        = (*Update)(nil)
	_ builder.CommandWithExecute     = (*Update)(nil)
	_ builder.CommandWithLogger      = (*Update)(nil)
	_ builder.CommandWithArgs        = (*Update)(nil)
	_ builder.CommandWithFlags       = (*Update)(nil)
)

func (*Update) Usage() string {
	return `update nameOrUUID --data '{"key": "any new data"}`
}

func (*Update) Docs() builder.Docs {
	return builder.Docs{
		Short: "Update a Turbine Secret",
		Long:  `This command will update the specified Turbine Secret's data.`,
		Example: `meroxa secrets update nameOrUUID --data '{"key": "value"}' 
		or 
		meroxa secrets update nameOrUUID `,
	}
}

func (d *Update) Config(cfg config.Config) {
	d.config = cfg
}

func (d *Update) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

// ParseArgs implements builder.CommandWithArgs.
func (d *Update) ParseArgs(args []string) error {
	if len(args) > 0 {
		d.args.nameOrUUID = args[0]
	}
	return nil
}

func (d *Update) BasicClient(client global.BasicClient) {
	d.client = client
}

func (d *Update) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Update) Execute(ctx context.Context) error {
	var err error
	getSecrets, err := RetrieveSecretsID(ctx, d.client, d.args.nameOrUUID)
	if err != nil {
		return err
	}

	fmt.Println("To proceed, enter new data for each key or press enter to skip. ")
	for k := range getSecrets.Items[0].Data {
		fmt.Printf("\n %q: ", k)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n') //nolint:shadow
		if err != nil {
			return err
		}
		if len(strings.TrimRight(input, "\r\n")) != 0 {
			getSecrets.Items[0].Data[k] = strings.TrimRight(input, "\r\n")
		}
	}

	if d.flags.Data != "" {
		appendData := make(map[string]interface{})
		err := json.Unmarshal([]byte(d.flags.Data), &appendData) //nolint:shadow
		if err != nil {
			return err
		}
		for key, string := range appendData {
			getSecrets.Items[0].Data[key] = string
		}
	}

	d.logger.Infof(ctx, "Updating secret %q...", d.args.nameOrUUID)
	response, err := d.client.CollectionRequest(ctx, "PATCH", collectionName, getSecrets.Items[0].ID, getSecrets.Items[0], nil)
	if err != nil {
		return err
	}

	updatedSecret := &Secrets{}
	err = json.NewDecoder(response.Body).Decode(&updatedSecret)
	if err != nil {
		return err
	}

	d.logger.Infof(ctx, "Secret %q has been updated.", updatedSecret.Name)
	d.logger.JSON(ctx, updatedSecret)

	return nil
}
