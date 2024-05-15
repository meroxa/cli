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

	"github.com/google/uuid"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

type Deploy struct {
	flags struct {
		Path string `long:"path" usage:"Path to the app directory (default is local directory)"`
	}

	client global.BasicClient
	logger log.Logger
	path   string
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

func switchToAppDirectory(appPath string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return pwd, err
	}
	return pwd, os.Chdir(appPath)
}

func shouldSkipDir(fi os.FileInfo) bool {
	if !fi.IsDir() {
		return false
	}

	switch fi.Name() {
	case ".git", "fixtures", "node_modules":
		return true
	}

	return false
}

func (d *Deploy) gzipConduitApp(src string, buf io.Writer) error {
	// Grab the directory we care about (app's directory)
	appDir := filepath.Base(src)

	// Change to parent's app directory
	pwd, err := switchToAppDirectory(filepath.Dir(src))
	if err != nil {
		return err
	}

	zipWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(zipWriter)

	err = filepath.Walk(appDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if shouldSkipDir(fi) {
			return filepath.SkipDir
		}
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(file)
		if err := tarWriter.WriteHeader(header); err != nil { //nolint:govet
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		if !fi.IsDir() {
			var data *os.File
			data, err = os.Open(file)
			defer func(data *os.File) {
				err = data.Close()
				if err != nil {
					panic(err.Error())
				}
			}(data)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := tarWriter.Close(); err != nil {
		return err
	}
	if err := zipWriter.Close(); err != nil {
		return err
	}

	return os.Chdir(pwd)
}

func (d *Deploy) Execute(ctx context.Context) error {
	var err error

	d.logger.StartSpinner("\t", "Starting app deploy, zipping and uploading application tar.")

	d.path, err = GetPath(d.flags.Path)
	if err != nil {
		return fmt.Errorf("error getting conduit app path - %s", err)
	}

	var buf bytes.Buffer
	err = d.gzipConduitApp(d.path, &buf)
	if err != nil {
		return fmt.Errorf("error zipping conduit app repository - %s", err)
	}

	dFile := fmt.Sprintf("conduit-%s.tar.gz", uuid.NewString())
	fileToWrite, err := os.OpenFile(dFile, os.O_CREATE|os.O_RDWR, os.FileMode(0o777)) //nolint:gomnd
	defer func(fileToWrite *os.File) {
		if err = fileToWrite.Close(); err != nil {
			panic(err.Error())
		}

		d.logger.StartSpinner("\t", fmt.Sprintf("Removing %q...", dFile))
		if err = os.Remove(dFile); err != nil {
			d.logger.StopSpinnerWithStatus(fmt.Sprintf("\t Something went wrong trying to remove %q", dFile), log.Failed)
		} else {
			d.logger.StopSpinnerWithStatus(fmt.Sprintf("Removed %q", dFile), log.Successful)
		}
	}(fileToWrite)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fileToWrite, &buf); err != nil {
		return err
	}

	files := map[string]string{
		"archive": fileToWrite.Name(),
	}

	response, err := d.client.CollectionRequestMultipart(
		ctx,
		http.MethodPost,
		deploymentCollection,
		"",
		nil,
		nil,
		files,
	)
	if err != nil {
		return err
	}

	var j map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&j)

	if j["state"] == "failed" {
		return fmt.Errorf("error deploying application, application state - %s", j["state"])
	}

	d.logger.StopSpinnerWithStatus("Application successfully deployed!", log.Successful)

	return nil
}
