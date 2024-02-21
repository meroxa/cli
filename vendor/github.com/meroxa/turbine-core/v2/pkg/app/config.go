package app

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/meroxa/turbine-core/v2/pkg/ir"
)

type Config struct {
	Name     string            `json:"name"`
	Fixtures map[string]string `json:"fixtures"`
	Language ir.Lang           `json:"language"`
}

// validateConfig will check if app.json contains information required.
func (c *Config) validateConfig() error {
	if c.Name == "" {
		return errors.New("application name is required to be specified in your app.json")
	}
	return nil
}

var ReadConfig = func(appName, appPath string) (Config, error) {
	if appPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatalf("unable to locate executable path; error: %s", err)
		}
		appPath = path.Dir(exePath)
	}

	b, err := os.ReadFile(filepath.Join(appPath, "app.json"))
	if err != nil {
		return Config{}, err
	}

	var ac Config
	err = json.Unmarshal(b, &ac)
	if err != nil {
		return Config{}, err
	}

	if appName != "" {
		ac.Name = appName
	}
	err = ac.validateConfig()
	if err != nil {
		return Config{}, err
	}

	return ac, nil
}
