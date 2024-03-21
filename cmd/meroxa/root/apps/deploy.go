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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
)

type Deploy struct {
	flags struct {
		Path string `long:"path" usage:"Path to the app directory (default is local directory)"`
		Spec string `long:"spec" usage:"Deployment specification version to use to build and deploy the app" hidden:"true"`
	}

	client global.BasicClient
	config config.Config
	logger log.Logger

	appName string
}

var (
	_ builder.CommandWithBasicClient = (*Deploy)(nil)
	_ builder.CommandWithConfig      = (*Deploy)(nil)
	_ builder.CommandWithDocs        = (*Deploy)(nil)
	_ builder.CommandWithExecute     = (*Deploy)(nil)
	_ builder.CommandWithFlags       = (*Deploy)(nil)
	_ builder.CommandWithLogger      = (*Deploy)(nil)
)

func (*Deploy) Usage() string {
	return "deploy [--path pwd]"
}

func (*Deploy) Docs() builder.Docs {
	return builder.Docs{
		Short: "Deploy a Turbine Data Application",
		Long: `This command will deploy the application specified in '--path'
(or current working directory if not specified) to our Meroxa Platform.
If deployment was successful, you should expect an application you'll be able to fully manage
`,
		Example: `meroxa apps deploy # assumes you run it from the app directory
meroxa apps deploy --path ./my-app
`,
	}
}

func (d *Deploy) Config(cfg config.Config) {
	d.config = cfg
}

func (d *Deploy) BasicClient(client global.BasicClient) {
	d.client = client

	// deployments needs to ensure enough time to complete
	if !global.ClientWithCustomTimeout() {
		d.client.SetTimeout(60 * time.Second)
	}
}

func (d *Deploy) Flags() []builder.Flag {
	return builder.BuildFlags(&d.flags)
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Deploy) Execute(ctx context.Context) error {
	var err error

	// TODO : add conduit logic for deploy

	response, err := d.client.CollectionRequestMultipart(
		ctx,
		http.MethodPost,
		collectionName,
		"",
		nil,
		nil,
		map[string]string{}, // TODO: change back to files from above
	)
	if err != nil {
		return err
	}

	apps := &Application{}
	err = json.NewDecoder(response.Body).Decode(&apps)
	if err != nil {
		return err
	}

	dashboardURL := fmt.Sprintf("%s/apps/%s/detail", global.GetMeroxaTenantURL(), apps.ID)
	output := fmt.Sprintf("Application %q successfully deployed!\n\n  ✨ To view your application, visit %s",
		d.appName, dashboardURL)

	d.logger.StopSpinnerWithStatus(output, log.Successful)

	return nil
}
