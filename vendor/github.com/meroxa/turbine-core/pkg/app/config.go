package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/meroxa/turbine-core/pkg/ir"
)

type Config struct {
	Name      string            `json:"name"`
	Pipeline  string            `json:"pipeline"` // TODO: Eventually remove support for providing a pipeline if we need to
	Resources map[string]string `json:"resources"`
	Language  ir.Lang           `json:"language"`
}

// validateAppConfig will check if app.json contains information required
func (c *Config) validateConfig() error {
	if c.Name == "" {
		return errors.New("application name is required to be specified in your app.json")
	}
	return nil
}

// setPipelineName will check if Pipeline was specified via app.json
// otherwise, pipeline name will be set with the format of `turbine-pipeline-{Name}`
func (c *Config) setPipelineName() {
	if c.Pipeline == "" {
		c.Pipeline = fmt.Sprintf("turbine-pipeline-%s", c.Name)
	}
}

var ReadConfig = func(appName, appPath string) (Config, error) {
	if appPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatalf("unable to locate executable path; error: %s", err)
		}
		appPath = path.Dir(exePath)
	}

	b, err := os.ReadFile(appPath + "/" + "app.json")
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

	ac.setPipelineName()
	return ac, nil
}
