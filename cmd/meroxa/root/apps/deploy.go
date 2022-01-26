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
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type Deploy struct {
	logger log.Logger
}

var (
	_ builder.CommandWithDocs = (*Deploy)(nil)
)

func (d *Deploy) Usage() string {
	return "deploy PATH"
}

func (*Deploy) Docs() builder.Docs {
	return builder.Docs{
		Short: "Deploy the current Meroxa Data Application",
	}
}

func (d *Deploy) Execute(ctx context.Context) error {
	// TODO:
	// - Generate wrapped binary (main)
	// - Deploy Docker Image (Functions)
	// -- Build docker image
	// -- Push container image to registry
	// - Run application (build pipeline)
	return nil
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}
