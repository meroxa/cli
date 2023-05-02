package deploy

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/meroxa/turbine-core/pkg/ir"

	"github.com/meroxa/turbine-go"
)

//go:embed template/*
var templateFS embed.FS

// TurbineDockerfileTrait will be used to replace data evaluations to generate an according Dockerfile
type TurbineDockerfileTrait struct {
	AppName   string
	GoVersion string
}

// CreateDockerfile will be used from the CLI to generate a new Dockerfile based on the app image
func CreateDockerfile(appName, pwd, specVersion string) error {
	if appName == "" {
		ac, err := turbine.ReadAppConfig(appName, pwd)
		if err != nil {
			log.Fatalln(err)
		}
		appName = ac.Name
	}

	fileName := "Dockerfile"

	// TODO: Remove this once 0.2.0 is default for Go
	if specVersion != ir.SpecVersion_0_2_0 {
		fileName = "Dockerfile.old"
	}

	t, err := template.ParseFS(templateFS, filepath.Join("template", fileName))
	if err != nil {
		return err
	}

	dockerfile := TurbineDockerfileTrait{
		AppName:   appName,
		GoVersion: "1.20",
	}

	f, err := os.Create(filepath.Join(pwd, "Dockerfile"))
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
