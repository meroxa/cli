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

package deploy

import (
	"embed"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

const (
	templateDir = "template"
	goVersion   = "1.17"
)

//go:embed template/*
var templateFS embed.FS

// TurbineDockerfileTrait will be used to replace data evaluations to generate an according Dockerfile
type TurbineDockerfileTrait struct {
	AppName   string
	GoVersion string
}

// CreateDockerfile will be used from the CLI to generate a new Dockerfile based on the app image
func CreateDockerfile(pwd string) error {
	fileName := "Dockerfile"
	appName := path.Base(pwd)
	t, err := template.ParseFS(templateFS, filepath.Join(templateDir, fileName))
	if err != nil {
		return err
	}

	dockerfile := TurbineDockerfileTrait{
		AppName:   appName,
		GoVersion: goVersion,
	}

	f, err := os.Create(filepath.Join(pwd, fileName))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	err = t.Execute(f, dockerfile)
	if err != nil {
		return err
	}
	return nil
}
