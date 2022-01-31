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
	"path"
)

type Deploy struct {
	args struct {
		Path string
	}
	logger log.Logger
}

var (
	_ builder.CommandWithDocs = (*Deploy)(nil)
)

func (d *Deploy) Usage() string {
	return "deploy [PATH]"
}

func (*Deploy) Docs() builder.Docs {
	return builder.Docs{
		Short: "Deploy the current Meroxa Data Application",
	}
}

func (d *Deploy) Execute(ctx context.Context) error {
	var projPath string
	if p := d.args.Path; p != "" {
		projPath = p
	} else {
		projPath = "."
	}
	projName := path.Base(projPath)

	fqImageName := prependAccount(projName)

	// build image
	err := buildImage(ctx, d.logger, ".", fqImageName)
	if err != nil {
		d.logger.Errorf(ctx, "unable to build image; %s", err)
	}

	// push image
	err = pushImage(ctx, d.logger, fqImageName)
	if err != nil {
		d.logger.Errorf(ctx, "unable to push image; %s", err)
	}

	// build go app
	err = buildGoApp(ctx, d.logger, projPath, true)
	if err != nil {
		return err
	}

	// deploy data app
	err = deployApp(ctx, d.logger, projPath, projName, fqImageName)
	if err != nil {
		d.logger.Errorf(ctx, "unable to deploy app; %s", err)
	}

	return nil
}

func (d *Deploy) Logger(logger log.Logger) {
	d.logger = logger
}

func (d *Deploy) ParseArgs(args []string) error {
	if len(args) > 0 {
		d.args.Path = args[0]
	}

	return nil
}
