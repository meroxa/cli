package turbinecli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/meroxa/cli/log"
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

// GitChecks prints warnings about uncommitted tracked and untracked files.
func GitChecks(ctx context.Context, l log.Logger, path string) error {
	// temporarily switching to the app's directory
	pwd, err := switchToAppDirectory(path)

	cmd := exec.Command("git", "status", "--porcelain=v2")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	all := string(output)
	lines := strings.Split(strings.TrimSpace(all), "\n")
	if len(lines) > 0 && lines[0] != "" {
		cmd = exec.Command("git", "status")
		output, err = cmd.Output()
		if err != nil {
			return err
		}
		l.Error(ctx, string(output))
		err = os.Chdir(pwd)
		if err != nil {
			return err
		}
		return fmt.Errorf("unable to proceed with deployment because of uncommitted changes")
	}
	return os.Chdir(pwd)
}

// GetPipelineUUID parses the deploy output when it was successful to determine the pipeline UUID to create.
func GetPipelineUUID(output string) string {
	// Example output:
	// 2022/03/16 13:21:36 pipeline created: "turbine-pipeline-simple" ("049760a8-a3d2-44d9-b326-0614c09a3f3e").
	re := regexp.MustCompile(`pipeline created:."[a-zA-Z]+-[a-zA-Z]+-[a-zA-Z]+".(\([^)]*\))`)
	res := re.FindStringSubmatch(output)[1]
	res = strings.Trim(res, "()\"")
	return res
}

// ValidateBranch validates the deployment is being performed from one of the allowed branches
func ValidateBranch(path string) error {
	// temporarily switching to the app's directory
	pwd, err := switchToAppDirectory(path)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	branchName := strings.TrimSpace(string(output))
	if branchName != "main" && branchName != "master" {
		return fmt.Errorf("deployment allowed only from 'main' or 'master' branch, not %s", branchName)
	}

	err = os.Chdir(pwd)
	if err != nil {
		return err
	}
	return nil
}

// GetGitSha will return the latest gitSha that will be used to create an application
func GetGitSha(path string) (string, error) {
	// temporarily switching to the app's directory
	pwd, err := switchToAppDirectory(path)
	if err != nil {
		return "", err
	}

	// Gets latest git sha
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	err = os.Chdir(pwd)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// switchToAppDirectory switches temporarily to the application's directory
func switchToAppDirectory(path string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return pwd, err
	}
	return pwd, os.Chdir(path)
}
