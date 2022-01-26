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
	"os"
	"os/exec"
	"path"
	"strings"
)

type Run struct {
	args struct {
		Path string
	}

	logger log.Logger
}

var (
	_ builder.CommandWithDocs = (*Run)(nil)
)

func (r *Run) Usage() string {
	return "run [PATH]"
}

func (*Run) Docs() builder.Docs {
	return builder.Docs{
		Short: "execute a Meroxa Data Application locally",
	}
}

func (r *Run) Execute(ctx context.Context) error {
	var projPath string
	if p := r.args.Path; p != "" {
		projPath = p
	} else {
		projPath = "."
	}
	buildPath := strings.Join([]string{projPath, "..."}, "/")
	r.logger.Info(ctx, "building apps...\n")
	cmd := exec.Command("go", "build", buildPath)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// TODO: parse output for build errors
	r.logger.Info(ctx, "build complete!")

	// apps name
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	projName := path.Base(pwd)

	cmd = exec.Command("./" + projName)
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	r.logger.Infof(ctx, "Running %s:", projName)
	r.logger.Info(ctx, string(stdout))

	return nil
}

func (r *Run) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Run) ParseArgs(args []string) error {
	if len(args) > 0 {
		r.args.Path = args[0]
	}

	return nil
}
