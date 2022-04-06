package deploy

import (
	"embed"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"
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
	t, err := template.ParseFS(templateFS, filepath.Join("template", fileName))
	if err != nil {
		return err
	}

	dockerfile := TurbineDockerfileTrait{
		AppName:   appName,
		GoVersion: "1.17",
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
