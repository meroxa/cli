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
	"context"
	"fmt"

	"github.com/skratchdot/open-golang/open"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

var (
	_ builder.CommandWithDocs        = (*Open)(nil)
	_ builder.CommandWithLogger      = (*Open)(nil)
	_ builder.CommandWithExecute     = (*Open)(nil)
	_ builder.CommandWithArgs        = (*Open)(nil)
	_ builder.CommandWithFlags       = (*Open)(nil)
	_ builder.CommandWithBasicClient = (*List)(nil)
)

type Opener interface {
	Start(string) error
}

type browserOpener struct{}

func (b *browserOpener) Start(URL string) error {
	return open.Start(URL)
}

func (o *Open) BasicClient(client global.BasicClient) {
	o.client = client
}

type Open struct {
	Opener

	client global.BasicClient

	logger log.Logger
	// path   string

	args struct {
		NameOrUUID string
	}
	flags struct {
		Path string `long:"path" usage:"Path to the app directory (default is local directory)"`
	}
}

func (o *Open) Usage() string {
	return "open [--path pwd]"
}

func (o *Open) Flags() []builder.Flag {
	return builder.BuildFlags(&o.flags)
}

func (o *Open) ParseArgs(args []string) error {
	if len(args) > 0 {
		o.args.NameOrUUID = args[0]
	}

	return nil
}

func (o *Open) Docs() builder.Docs {
	return builder.Docs{
		Short: "Open the link to a Conduit Data Application in the Dashboard",
		Example: `meroxa apps open # assumes that the Application is in the current directory
meroxa apps open NAMEorUUID`,
	}
}

func (o *Open) Execute(ctx context.Context) error {
	if o.Opener == nil {
		o.Opener = &browserOpener{}
	}

	apps, err := RetrieveApplicationByNameOrID(ctx, o.client, o.args.NameOrUUID)
	if err != nil {
		return err
	}

	// open a browser window to the application details
	dashboardURL := fmt.Sprintf("%s/apps/%s/detail", global.GetMeroxaTenantURL(), apps.Items[0].ID)
	err = o.Start(dashboardURL)
	if err != nil {
		o.logger.Errorf(ctx, "can't open browser to URL %s\n", dashboardURL)
	}
	return err
}

func (o *Open) Logger(logger log.Logger) {
	o.logger = logger
}
