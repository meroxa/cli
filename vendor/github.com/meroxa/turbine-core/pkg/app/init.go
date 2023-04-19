package app

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/meroxa/turbine-core/pkg/ir"
)

type AppInit struct {
	AppName  string
	Language ir.Lang
	Path     string
}

// AppInitTemplate will be used to replace data evaluations provided by the user
type AppInitTemplate struct {
	AppName string
}

//go:embed templates/*
var templateFS embed.FS

func (a *AppInit) createAppDirectory() error {
	return os.MkdirAll(filepath.Join(a.Path, a.AppName), 0o755)
}

// createFixtures will create exclusively a fixtures folder and its content
func (a *AppInit) createFixtures() error {
	directory := "fixtures"

	err := os.Mkdir(filepath.Join(a.Path, a.AppName, directory), 0o755)
	if err != nil {
		return err
	}

	content, err := templateFS.ReadDir(filepath.Join("templates", string(a.Language), directory))
	if err != nil {
		return err
	}

	for _, f := range content {
		if !f.IsDir() && strings.Contains(f.Name(), "json") {
			if err = a.duplicateFile("fixtures/" + f.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

// duplicateFile reads from a template and write to a file located to a path provided by the user
func (a *AppInit) duplicateFile(fileName string) error {
	t, err := template.ParseFS(templateFS, filepath.Join("templates", string(a.Language), fileName))
	if err != nil {
		return err
	}

	appTrait := AppInitTemplate{
		AppName: a.AppName,
	}
	f, err := os.Create(filepath.Join(a.Path, a.AppName, fileName))
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
func (a *AppInit) listTemplateContent() ([]string, []string, error) {
	var files, directories []string

	content, err := templateFS.ReadDir(filepath.Join("templates", string(a.Language)))
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

func NewAppInit(appName string, language ir.Lang, path string) *AppInit {
	return &AppInit{
		AppName:  appName,
		Language: language,
		Path:     path,
	}
}

// Init will be used from the CLI to generate a new application directory based on the existing
// content on `/templates`.
func (a *AppInit) Init() error {
	err := a.createAppDirectory()
	if err != nil {
		return err
	}

	files, _, err := a.listTemplateContent()
	if err != nil {
		return err
	}

	for _, f := range files {
		err := a.duplicateFile(f)
		if err != nil {
			return err
		}
	}

	// TODO: Maybe write a recursive function that could take care of nested directories like this one
	err = a.createFixtures()
	if err != nil {
		return err
	}
	return nil
}
