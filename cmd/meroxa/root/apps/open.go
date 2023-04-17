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
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

var (
	_ builder.CommandWithDocs    = (*Open)(nil)
	_ builder.CommandWithLogger  = (*Open)(nil)
	_ builder.CommandWithExecute = (*Open)(nil)
	_ builder.CommandWithArgs    = (*Open)(nil)
	_ builder.CommandWithFlags   = (*Open)(nil)
)

type Open struct {
	logger log.Logger
	path   string

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
		Short: "Open the link to a Turbine Data Application in the Dashboard",
		Example: `meroxa apps open # assumes that the Application is in the current directory
meroxa apps open --path /my/app
meroxa apps open NAMEorUUID`,
	}
}

func (o *Open) Execute(ctx context.Context) error {
	nameOrUUID := o.args.NameOrUUID
	if nameOrUUID != "" && o.flags.Path != "" {
		return fmt.Errorf("supply either NameOrUUID argument or --path flag")
	}

	if nameOrUUID == "" {
		var err error
		if o.path, err = turbine.GetPath(o.flags.Path); err != nil {
			return err
		}

		config, err := turbine.ReadAppConfigFile(o.path)
		if err != nil {
			return err
		}
		nameOrUUID = config.Name
	}

	// open a browser window to the application details
	dashboardURL := fmt.Sprintf("https://dashboard.meroxa.io/apps/%s/detail", nameOrUUID)
	err := open.Start(dashboardURL)
	if err != nil {
		o.logger.Errorf(ctx, "can't open browser to URL %s\n", dashboardURL)
	}
	return err
}

func (o *Open) Logger(logger log.Logger) {
	o.logger = logger
}
