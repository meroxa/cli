package app

import (
	"embed"
	"github.com/meroxa/turbine-core/pkg/ir"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"
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

//go:embed all:templates
var templateFS embed.FS

func NewAppInit(appName string, language ir.Lang, path string) *AppInit {
	return &AppInit{
		AppName:  appName,
		Language: language,
		Path:     path,
	}
}

func (a *AppInit) applytemplate(srcDir, destDir, fileName string) error {
	t, err := template.ParseFS(templateFS, path.Join(srcDir, fileName))
	if err != nil {
		return err
	}

	appTrait := AppInitTemplate{
		AppName: a.AppName,
	}

	f, err := os.Create(filepath.Join(destDir, fileName))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	return t.Execute(f, appTrait)
}

// copyFile simply copies the file from srcDir to destDir (without applying a template)
func (a *AppInit) copyFile(srcDir, destDir, fileName string) error {
	srcPath := filepath.Join(srcDir, fileName)
	destPath := filepath.Join(destDir, fileName)

	srcFile, err := templateFS.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func (a *AppInit) duplicateFileInPath(srcDir, destDir, fileName string) error {
	notTemplateExtensions := []string{".jar"}

	fExtension := filepath.Ext(fileName)
	for _, ext := range notTemplateExtensions {
		if fExtension == ext {
			return a.copyFile(srcDir, destDir, fileName)
		}
	}
	return a.applytemplate(srcDir, destDir, fileName)
}

// listTemplateContentFromPath is used to return existing files and directories on a given path
func (a *AppInit) listTemplateContentFromPath(srcPath string) ([]string, []string, error) {
	var files, directories []string

	content, err := templateFS.ReadDir(srcPath)
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

func (a *AppInit) duplicateDirectory(srcDir, destDir string) error {
	// Create source directory
	err := os.MkdirAll(destDir, 0o755)
	if err != nil {
		return err
	}

	files, directories, err := a.listTemplateContentFromPath(srcDir)

	for _, fileName := range files {
		err = a.duplicateFileInPath(srcDir, destDir, fileName)
		if err != nil {
			return err
		}
	}

	for _, d := range directories {
		subSrcDir := filepath.Join(srcDir, d)
		subDestDir := filepath.Join(destDir, d)
		err = a.duplicateDirectory(subSrcDir, subDestDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// Init will be used from the CLI to generate a new application directory based on the existing
// content on `/templates`.
func (a *AppInit) Init() error {
	rootSrcDir := filepath.Join("templates", string(a.Language))
	rootDestDir := filepath.Join(a.Path, a.AppName)

	return a.duplicateDirectory(rootSrcDir, rootDestDir)
}
