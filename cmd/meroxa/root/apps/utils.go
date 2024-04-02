package apps

import (
	"os"
	"path/filepath"
)

type AppConfig struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Resources   map[string]string `json:"resources"`
	Vendor      string            `json:"vendor"`
	ModuleInit  string            `json:"module_init"`
}

var prefetched *AppConfig

func GetPath(flag string) (string, error) {
	if flag == "" {
		flag = "."
	}
	var err error
	flag, err = filepath.Abs(flag)
	if err != nil {
		return "", err
	}
	return flag, nil
}

// SwitchToAppDirectory switches temporarily to the application's directory.
func SwitchToAppDirectory(appPath string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return pwd, err
	}
	return pwd, os.Chdir(appPath)
}
