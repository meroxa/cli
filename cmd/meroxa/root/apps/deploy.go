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
	"io"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
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
	// TODO:
	// - Generate wrapped binary (main)
	// - Deploy Docker Image (Functions)
	// -- Build docker image
	// -- Push container image to registry
	// - Run application (build pipeline)
	var projPath string
	if p := d.args.Path; p != "" {
		projPath = p
		//os.Chdir(projPath)
	} else {
		projPath = "."
	}
	projName := path.Base(projPath)

	err := d.buildImage(ctx, ".", projName)
	if err != nil {
		d.logger.Errorf(ctx, "unable to build image; %s", err)
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

func (d *Deploy) buildImage(ctx context.Context, path string, name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		d.logger.Errorf(ctx, "unable to init docker client; %s", err)
	}
	// Read local Dockerfile
	tar, err := archive.TarWithOptions(".", &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{"simple", ".git", "fixtures"},
	})
	if err != nil {
		d.logger.Errorf(ctx, "unable to create tar; %s", err)
	}

	buildOptions := types.ImageBuildOptions{
		Context:    tar,
		Dockerfile: "Dockerfile",
		Remove:     true,
		Tags:       []string{imageName(name)}}

	resp, err := cli.ImageBuild(
		ctx,
		tar,
		buildOptions,
	)
	if err != nil {
		d.logger.Errorf(ctx, "unable to build docker image; %s", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		d.logger.Errorf(ctx, "unable to read image build response; %s", err)
	}
	return nil
}

func imageName(name string) string {
	scope := os.Getenv("DOCKER_HUB_USERNAME")
	return strings.Join([]string{scope, name}, "/")
}
