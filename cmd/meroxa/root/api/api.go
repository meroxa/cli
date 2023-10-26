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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
)

var (
	_ builder.CommandWithDocs        = (*API)(nil)
	_ builder.CommandWithArgs        = (*API)(nil)
	_ builder.CommandWithBasicClient = (*API)(nil)
	_ builder.CommandWithLogger      = (*API)(nil)
	_ builder.CommandWithExecute     = (*API)(nil)
)

type API struct {
	client global.BasicClient
	logger log.Logger

	args struct {
		Method string
		Path   string
		Body   interface{}
		ID     string
	}
}

func (a *API) Usage() string {
	return "api METHOD PATH [body]"
}

func (a *API) Docs() builder.Docs {
	return builder.Docs{
		Short: "Invoke Meroxa API",
		Example: `
meroxa api GET {collection} {id}
meroxa api POST {collection} {id} '{"type":"postgres", "name":"pg", "url":"postgres://.."}'`,
	}
}

func (a *API) BasicClient(client global.BasicClient) {
	a.client = client
}

func (a *API) Logger(logger log.Logger) {
	a.logger = logger
}

func (a *API) ParseArgs(args []string) error {
	if len(args) < 2 {
		return errors.New("requires METHOD and PATH")
	}

	a.args.Method = strings.ToUpper(args[0]) // so either a method such as `get` or `GET` works
	a.args.Path = args[1]

	if len(args) > 2 {
		a.args.Body = args[2]
	}

	return nil
}

func (a *API) Execute(ctx context.Context) error {

	resp, err := a.client.URLRequest(ctx, a.args.Method, a.args.Path, a.args.Body, nil, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var prettyJSON bytes.Buffer

	if err = json.Indent(&prettyJSON, b, "", "\t"); err != nil {
		prettyJSON.Write(b)
	}

	a.logger.Infof(ctx, "> %s %s", a.args.Method, a.args.Path)
	a.logger.Infof(ctx, "< %s %s", resp.Status, resp.Proto)
	for k, v := range resp.Header {
		a.logger.Infof(ctx, "< %s %s", k, strings.Join(v, " "))
	}

	a.logger.Info(ctx, prettyJSON.String())
	a.logger.JSON(ctx, prettyJSON.String())

	return nil
}
