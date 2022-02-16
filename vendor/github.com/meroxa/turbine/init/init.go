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

package init

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	templateDir  = "template"
	pipelineName = "default"
)

//go:embed template/*
var templateFS embed.FS

// TurbineAppTrait will be used to replace data evaluations provided by the user
type TurbineAppTrait struct {
	Name     string
	Pipeline string
}

// createAppDirectory is where new files will be created. It'll be named as the application name
func createAppDirectory(path, appName string) error {
	return os.MkdirAll(filepath.Join(path, appName), 0755)
}

// createFixtures will create exclusively a fixtures folder and its content
func createFixtures(path, appName string) error {
	directory := "fixtures"
	fileName := "README.md"

	err := os.Mkdir(filepath.Join(path, appName, directory), 0755)
	if err != nil {
		return err
	}

	t, err := template.ParseFS(templateFS, filepath.Join(templateDir, directory, fileName))
	if err != nil {
		return err
	}

	appJSON := TurbineAppTrait{
		Name:     appName,
		Pipeline: pipelineName,
	}

	f, err := os.Create(filepath.Join(path, appName, directory, fileName))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	err = t.Execute(f, appJSON)
	if err != nil {
		return err
	}
	return nil
}

// duplicateFile reads from a template and write to a file located to a path provided by the user
func duplicateFile(fileName, path, appName string) error {
	t, err := template.ParseFS(templateFS, filepath.Join(templateDir, fileName))
	if err != nil {
		return err
	}

	appTrait := TurbineAppTrait{
		Name:     appName,
		Pipeline: pipelineName, // this could be provided by the user
	}

	f, err := os.Create(filepath.Join(path, appName, fileName))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	err = t.Execute(f, appTrait)
	if err != nil {
		return err
	}
	return nil
}

// listTemplateContent is used to return existing files and directories on a given path
func listTemplateContent() ([]string, []string, error) {
	var files, directories []string

	content, err := templateFS.ReadDir(templateDir)
	if err != nil {
		return files, directories, err
	}

	for _, f := range content {
		if f.IsDir() {
			directories = append(directories, f.Name())
		} else {
			files = append(files, f.Name())
		}
	}
	return files, directories, nil
}

// Init will be used from the CLI to generate a new application directory based on the existing
// content on `/template`.
func Init(path, appName string) error {
	err := createAppDirectory(path, appName)
	if err != nil {
		return err
	}

	files, _, err := listTemplateContent()
	if err != nil {
		return err
	}

	for _, f := range files {
		err := duplicateFile(f, path, appName)
		if err != nil {
			return err
		}
	}

	// TODO: Maybe write a recursive function that could take care of nested directories like this one
	err = createFixtures(path, appName)
	if err != nil {
		return err
	}
	return nil
}
