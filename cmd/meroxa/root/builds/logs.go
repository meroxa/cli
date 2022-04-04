/*
Copyright © 2021 Meroxa Inc

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

package builds

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Logs)(nil)
	_ builder.CommandWithArgs    = (*Logs)(nil)
	_ builder.CommandWithClient  = (*Logs)(nil)
	_ builder.CommandWithLogger  = (*Logs)(nil)
	_ builder.CommandWithExecute = (*Logs)(nil)
)

type buildLogsClient interface {
	GetBuildLogs(ctx context.Context, uuid string) (*http.Response, error)
}

type Logs struct {
	client buildLogsClient
	logger log.Logger

	args struct {
		UUID string
	}
}

func (d *Logs) Usage() string {
	return "logs [UUID]"
}

func (d *Logs) Docs() builder.Docs {
	return builder.Docs{
		Short: "List a Meroxa Process Build's Logs",
	}
}

func (d *Logs) Execute(ctx context.Context) error {
	response, err := d.client.GetBuildLogs(ctx, d.args.UUID)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	d.logger.Info(ctx, string(body))

	return nil
}

func (d *Logs) Client(client meroxa.Client) {
	d.client = client
}

func (d *Logs) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Logs) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires build UUID")
	}

	d.args.UUID = args[0]
	return nil
}
