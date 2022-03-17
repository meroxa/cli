package turbinecli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type AppConfig struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

func GetPath(pwd string) string {
	if pwd != "" {
		return pwd
	}
	return "."
}

// GetLang will return language defined either by `--lang` or the one defined by user in the app.json.
func GetLang(flag, pwd string) (string, error) {
	if flag != "" {
		return flag, nil
	}

	lang, err := GetLangFromAppJSON(pwd)
	if err != nil {
		return "", fmt.Errorf("flag --lang is required unless lang is specified in your app.json")
	}
	return lang, nil
}

// GetLangFromAppJSON returns specified language in users' app.json.
func GetLangFromAppJSON(pwd string) (string, error) {
	appConfig, err := readConfigFile(pwd)
	if err != nil {
		return "", err
	}

	if appConfig.Language == "" {
		return "", fmt.Errorf("`language` should be specified in your app.json")
	}
	return appConfig.Language, nil
}

// GetAppNameFromAppJSON returns specified app name in users' app.json.
func GetAppNameFromAppJSON(pwd string) (string, error) {
	appConfig, err := readConfigFile(pwd)
	if err != nil {
		return "", err
	}

	if appConfig.Name == "" {
		return "", fmt.Errorf("`name` should be specified in your app.json")
	}
	return appConfig.Name, nil
}

// readConfigFile will read the content of an app.json based on path.
func readConfigFile(appPath string) (AppConfig, error) {
	var appConfig AppConfig

	appConfigPath := path.Join(appPath, "app.json")
	appConfigBytes, err := os.ReadFile(appConfigPath)
	if err != nil {
		return appConfig, fmt.Errorf("%v\n"+
			"We couldn't find an app.json file on path %q. Maybe try in another using `--path`", err, appPath)
	}
	if err := json.Unmarshal(appConfigBytes, &appConfig); err != nil {
		return appConfig, err
	}

	return appConfig, nil
}
