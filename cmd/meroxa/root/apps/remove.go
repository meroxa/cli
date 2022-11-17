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
	"net/http"
	"os"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

type removeAppClient interface {
	DeleteApplicationEntities(ctx context.Context, name string) (*http.Response, error)
	AddHeader(key, value string)
}

type Remove struct {
	client     removeAppClient
	logger     log.Logger
	turbineCLI turbine.CLI
	path       string

	args struct {
		NameOrUUID string
	}
	flags struct {
		Path string `long:"path" usage:"Path to the app directory (default is local directory)"`
	}
}

func (r *Remove) Usage() string {
	return `remove [NameOrUUID] [--path pwd]`
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
meroxa apps remove NAME`,
		Beta: true,
	}
}

func (r *Remove) Execute(ctx context.Context) error {
	var turbineLibVersion string
	nameOrUUID := r.args.NameOrUUID
	if nameOrUUID != "" && r.flags.Path != "" {
		return fmt.Errorf("supply either NamrOrUUID argument or path flag")
	}

	if nameOrUUID == "" {
		var err error
		if r.path, err = turbine.GetPath(r.flags.Path); err != nil {
			return err
		}

		config, err := turbine.ReadConfigFile(r.path)
		if err != nil {
			return err
		}
		nameOrUUID = config.Name

		if r.turbineCLI == nil {
			r.turbineCLI, err = getTurbineCLIFromLanguage(r.logger, config.Language, r.path)
			if err != nil {
				return err
			}
		}

		if turbineLibVersion, err = r.turbineCLI.GetVersion(ctx); err != nil {
			return err
		}
		addTurbineHeaders(r.client, config.Language, turbineLibVersion)
	}

	if os.Getenv("UNIT_TEST") == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("To proceed, type %q or re-run this command with --force\n▸ ", nameOrUUID)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if nameOrUUID != strings.TrimSuffix(input, "\n") {
			return errors.New("action aborted")
		}
	}

	r.logger.Infof(ctx, "Removing application %q...", nameOrUUID)

	res, err := r.client.DeleteApplicationEntities(ctx, nameOrUUID)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	r.logger.Infof(ctx, "Application %q successfully removed", nameOrUUID)
	r.logger.JSON(ctx, res)

	return nil
}

func (r *Remove) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Remove) Client(client meroxa.Client) {
	r.client = client
}

func (r *Remove) ParseArgs(args []string) error {
	if len(args) > 0 {
		r.args.NameOrUUID = args[0]
	}

	return nil
}

func (r *Remove) Aliases() []string {
	return []string{"rm", "delete"}
}

var (
	_ builder.CommandWithDocs    = (*Remove)(nil)
	_ builder.CommandWithAliases = (*Remove)(nil)
	_ builder.CommandWithArgs    = (*Remove)(nil)
	_ builder.CommandWithFlags   = (*Remove)(nil)
	_ builder.CommandWithClient  = (*Remove)(nil)
	_ builder.CommandWithLogger  = (*Remove)(nil)
	_ builder.CommandWithExecute = (*Remove)(nil)
)
