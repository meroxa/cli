/*
Copyright © 2022 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apps

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

type Remove struct {
	client global.BasicClient
	logger log.Logger

	args struct {
		idOrName string
	}
	flags struct {
		Path  string `long:"path" usage:"Path to the app directory (default is local directory)"`
		Force bool   `long:"force" short:"f" default:"false" usage:"skip confirmation"`
	}
}

func (r *Remove) Usage() string {
	return `remove [ID or Name] [--path pwd]`
}

func (r *Remove) Flags() []builder.Flag {
	return builder.BuildFlags(&r.flags)
}

func (r *Remove) Docs() builder.Docs {
	return builder.Docs{
		Short: "Remove a Turbine Data Application",
		Long: `This command will remove the Application specified in '--path'
(or current working directory if not specified) previously deployed on our Meroxa Platform,
or the Application specified by the given name or UUID identifier.`,
		Example: `meroxa apps remove # assumes that the Application is in the current directory
meroxa apps remove --path /my/app
meroxa apps remove IDorName`,
	}
}

func (r *Remove) Execute(ctx context.Context) error {
	if !r.flags.Force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("To proceed, type %q or re-run this command with --force\n▸ ", r.args.idOrName)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if r.args.idOrName != strings.TrimRight(input, "\r\n") {
			return errors.New("action aborted")
		}
	}

	apps := &Applications{}
	apps, err := apps.RetrieveApplicationID(ctx, r.client, r.args.idOrName, r.flags.Path)
	if err != nil {
		return err
	}

	r.logger.Infof(ctx, "Removing application %q...", r.args.idOrName)
	response, err := r.client.CollectionRequest(ctx, "DELETE", collectionName, apps.Items[0].ID, nil, nil)
	if err != nil {
		return err
	}

	r.logger.Infof(ctx, "Application %q successfully removed", r.args.idOrName)
	r.logger.JSON(ctx, response)

	return nil
}

func (r *Remove) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Remove) BasicClient(client global.BasicClient) {
	r.client = client
}

func (r *Remove) ParseArgs(args []string) error {
	if len(args) > 0 {
		r.args.idOrName = args[0]
	}

	return nil
}

func (r *Remove) Aliases() []string {
	return []string{"rm", "delete"}
}

var (
	_ builder.CommandWithDocs        = (*Remove)(nil)
	_ builder.CommandWithAliases     = (*Remove)(nil)
	_ builder.CommandWithArgs        = (*Remove)(nil)
	_ builder.CommandWithFlags       = (*Remove)(nil)
	_ builder.CommandWithBasicClient = (*Remove)(nil)
	_ builder.CommandWithLogger      = (*Remove)(nil)
	_ builder.CommandWithExecute     = (*Remove)(nil)
)
