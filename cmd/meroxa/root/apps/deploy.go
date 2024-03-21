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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

type Deploy struct {
	flags struct {
		Path string `long:"path" usage:"Path to the app directory (default is local directory)"`
	}

	client  global.BasicClient
	logger  log.Logger
	appName string
	path    string
	gitSha  string
}

type ApplicationDeployment struct {
	ID                string                 `json:"id"`
	State             ApplicationState       `json:"state"`
	ApplicationSpec   map[string]interface{} `json:"app_spec"`
	ProcessorsPlugins map[string]interface{} `json:"processors_plugins"`
	PipelineFilenames map[string]interface{} `json:"pipelines_filenames"`
	Created           AppTime                `json:"created"`
	Updated           AppTime                `json:"updated"`
	Archive           string                 `json:"archive"`
}

var (
	_ builder.CommandWithBasicClient = (*Deploy)(nil)
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
		Short: "Deploy a Conduit Data Application",
		Long: `This command will deploy the application specified in '--path'
(or current working directory if not specified) to our Meroxa Platform.
If deployment was successful, you should expect an application you'll be able to fully manage
`,
		Example: `meroxa apps deploy # assumes you run it from the app directory
meroxa apps deploy --path ./my-app
`,
	}
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

func (d *Deploy) gzipConduitApp(src string, buf io.Writer) error {
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	// walk through every file in the folder
	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(file)

		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	if err := tw.Close(); err != nil {
		return err
	}
	if err := zr.Close(); err != nil {
		return err
	}
	//
	return nil
}

func (d *Deploy) getGitInfo(ctx context.Context) error {
	var err error
	if err = CheckUncommittedChanges(ctx, d.logger, d.path); err != nil {
		return err
	}

	d.gitSha, err = GetGitSha(ctx, d.path)
	return err
}

func (d *Deploy) Execute(ctx context.Context) error {
	var err error

	if err = d.getGitInfo(ctx); err != nil { //nolint:shadow
		return err
	}

	d.path, err = GetPath(d.flags.Path)
	if err != nil {
		return err
	}

	fmt.Println(d.path)

	var buf bytes.Buffer
	d.gzipConduitApp(d.path, &buf)

	// write the .tar.gzip
	file, err := os.CreateTemp("", "temp_conduit.tar.gzip")
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(file, &buf); err != nil {
		panic(err)
	}

	fmt.Println(file.Name())
	files := map[string]string{
		"archive": file.Name(),
	}

	fmt.Println(files)

	response, err := d.client.CollectionRequestMultipart(
		ctx,
		http.MethodPost,
		deploymentCollection,
		"",
		nil,
		nil,
		//		map[string]string{},
		files,
	)
	if err != nil {
		return err
	}

	apps := &Application{}
	err = json.NewDecoder(response.Body).Decode(&apps)
	if err != nil {
		return err
	}

	dashboardURL := fmt.Sprintf("%s/apps/%s/detail", global.GetMeroxaAPIURL(), apps.ID)
	fmt.Sprintf("Application %q successfully deployed!\n\n  ✨ To view your application, visit %s",
		d.appName, dashboardURL)

	//d.logger.StopSpinnerWithStatus(output, log.Successful)

	return nil
}
