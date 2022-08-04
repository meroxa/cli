package turbine

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
)

type AppConfig struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Pipeline    string            `json:"pipeline"` // TODO: Eventually remove support for providing a pipeline if we need to
	Resources   map[string]string `json:"resources"`
}

// validateAppConfig will check if app.json contains information required
func (c *AppConfig) validateAppConfig() error {
	if c.Name == "" {
		return errors.New("application name is required to be specified in your app.json")
	}
	return nil
}

// setPipelineName will check if Pipeline was specified via app.json
// otherwise, pipeline name will be set with the format of `turbine-pipeline-{Name}`
func (c *AppConfig) setPipelineName() {
	if c.Pipeline == "" {
		c.Pipeline = fmt.Sprintf("turbine-pipeline-%s", c.Name)
	}
}

var ReadAppConfig = func(appName, appPath string) (AppConfig, error) {
	if appPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatalf("unable to locate executable path; error: %s", err)
		}
		appPath = path.Dir(exePath)
	}

	b, err := os.ReadFile(appPath + "/" + "app.json")
	if err != nil {
		return AppConfig{}, err
	}

	var ac AppConfig
	err = json.Unmarshal(b, &ac)
	if err != nil {
		return AppConfig{}, err
	}

	if appName != "" {
		ac.Name = appName
	}
	err = ac.validateAppConfig()
	if err != nil {
		return AppConfig{}, err
	}

	ac.setPipelineName()
	return ac, nil
}
