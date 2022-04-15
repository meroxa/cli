package init

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	templateDir = "template"
)

//go:embed template/*
var templateFS embed.FS

// TurbineAppInitTrait will be used to replace data evaluations provided by the user
type TurbineAppInitTrait struct {
	AppName string
}

// createAppDirectory is where new files will be created. It'll be named as the application name
func createAppDirectory(path, appName string) error {
	return os.MkdirAll(filepath.Join(path, appName), 0755)
}

// createFixtures will create exclusively a fixtures folder and its content
func createFixtures(path, appName string) error {
	directory := "fixtures"

	err := os.Mkdir(filepath.Join(path, appName, directory), 0755)
	if err != nil {
		return err
	}

	content, err := templateFS.ReadDir(filepath.Join(templateDir, directory))
	if err != nil {
		return err
	}

	for _, f := range content {
		if !f.IsDir() && strings.Contains(f.Name(), "json") {
			if err = duplicateFile("fixtures/"+f.Name(), path, appName); err != nil {
				return err
			}
		}
	}

	return nil
}

// duplicateFile reads from a template and write to a file located to a path provided by the user
func duplicateFile(fileName, path, appName string) error {
	t, err := template.ParseFS(templateFS, filepath.Join(templateDir, fileName))
	if err != nil {
		return err
	}

	appTrait := TurbineAppInitTrait{
		AppName: appName,
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
func Init(appName, path string) error {
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
