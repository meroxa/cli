package turbine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/meroxa/cli/log"
)

const (
	GoLang     = "golang"
	JavaScript = "javascript"
	NodeJs     = "nodejs"
	Python     = "python"
	Python3    = "python3"
	Ruby       = "ruby"

	isTrue            = "true"
	AccountUUIDEnvVar = "MEROXA_ACCOUNT_UUID"

	IncompatibleTurbineVersion = `your Turbine library version is incompatible with the Meroxa CLI.
For guidance on updating to the latest version, visit:
https://docs.meroxa.com/beta-overview#updated-meroxa-cli-and-outdated-turbine-library`
)

type AppConfig struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Language    string            `json:"language"`
	Resources   map[string]string `json:"resources"`
	Vendor      string            `json:"vendor"`
	ModuleInit  string            `json:"module_init"`
}

var prefetched *AppConfig

type ApplicationResource struct {
	Name        string `json:"name"`
	Source      bool   `json:"source"`
	Destination bool   `json:"destination"`
	Collection  string `json:"collection"`
}

// GetResourceNamesFromString provides backward compatibility with turbine-go
// legacy resource listing format.
func GetResourceNamesFromString(s string) []ApplicationResource {
	resources := make([]ApplicationResource, 0)

	r := regexp.MustCompile(`\[(.+?)\]`)
	sliceString := r.FindStringSubmatch(s)
	if len(sliceString) == 0 {
		return resources
	}

	for _, n := range strings.Fields(sliceString[1]) {
		resources = append(resources, ApplicationResource{
			Name: n,
		})
	}

	return resources
}

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

// GetLangFromAppJSON returns specified language in users' app.json.
func GetLangFromAppJSON(ctx context.Context, l log.Logger, pwd string) (string, error) {
	l.StartSpinner("\t", "Determining application language from app.json...")
	appConfig, err := ReadConfigFile(pwd)
	if err != nil {
		l.StopSpinnerWithStatus("Something went wrong reading your app.json", log.Failed)
		return "", err
	}

	if appConfig.Language == "" {
		l.StopSpinnerWithStatus("`language` should be specified in your app.json", log.Failed)
		return "", fmt.Errorf("add key `language` to your app.json")
	}
	l.StopSpinnerWithStatus(fmt.Sprintf("Checked your language is %q", appConfig.Language), log.Successful)
	return appConfig.Language, nil
}

// GetAppNameFromAppJSON returns specified app name in users' app.json.
func GetAppNameFromAppJSON(ctx context.Context, l log.Logger, pwd string) (string, error) {
	l.StartSpinner("\t", "Reading application name from app.json...")
	appConfig, err := ReadConfigFile(pwd)
	if err != nil {
		return "", err
	}

	if appConfig.Name == "" {
		l.StopSpinnerWithStatus("`name` should be specified in your app.json", log.Failed)
		return "", fmt.Errorf("add `name` to your app.json")
	}
	l.StopSpinnerWithStatus(fmt.Sprintf("Checked your application name is %q", appConfig.Name), log.Successful)
	return appConfig.Name, nil
}

// SetModuleInitInAppJSON returns whether to run mod init.
func SetModuleInitInAppJSON(pwd string, skipInit bool) error {
	appConfig, err := ReadConfigFile(pwd)
	if err != nil {
		return err
	}
	appConfig.ModuleInit = "true"
	if skipInit {
		appConfig.ModuleInit = "false"
	}
	err = WriteConfigFile(pwd, appConfig)
	return err
}

// SetVendorInAppJSON returns whether to vendor modules.
func SetVendorInAppJSON(pwd string, vendor bool) error {
	appConfig, err := ReadConfigFile(pwd)
	if err != nil {
		return err
	}
	appConfig.Vendor = "false"
	if vendor {
		appConfig.Vendor = "true"
	}
	err = WriteConfigFile(pwd, appConfig)
	return err
}

// ReadConfigFile will read the content of an app.json based on path.
func ReadConfigFile(appPath string) (AppConfig, error) {
	var appConfig AppConfig

	if prefetched == nil || os.Getenv("UNIT_TEST") != "" {
		appConfigPath := path.Join(appPath, "app.json")
		appConfigBytes, err := os.ReadFile(appConfigPath)
		if err != nil {
			return appConfig, fmt.Errorf("could not find an app.json file on path %q."+
				" Try a different value for `--path`", appPath)
		}
		if err := json.Unmarshal(appConfigBytes, &appConfig); err != nil {
			return appConfig, err
		}
		prefetched = &appConfig
	}

	return *prefetched, nil
}

func WriteConfigFile(appPath string, cfg AppConfig) error {
	appConfigPath := path.Join(appPath, "app.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(appConfigPath, data, 0664) //nolint:gosec,gomnd
	if err != nil {
		return fmt.Errorf("%v\n"+
			"unable to update app.json file on path %q. Maybe try using a different value for `--path`", err, appPath)
	}
	return nil
}

func GitInit(ctx context.Context, appPath string) error {
	if appPath == "" {
		return errors.New("path is required")
	}

	isGitOlderThan228, err := checkGitVersion(ctx)
	if err != nil {
		return err
	}

	if !isGitOlderThan228 {
		cmd := exec.CommandContext(ctx, "git", "config", "--global", "init.defaultBranch", "main")
		if out, err := cmd.CombinedOutput(); err != nil {
			return errors.New(string(out))
		}
	}

	cmd := exec.CommandContext(ctx, "git", "init", appPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf(string(out))
	}

	if isGitOlderThan228 {
		cmd := exec.CommandContext(ctx, "git", "checkout", "-b", "main")
		cmd.Dir = appPath
		if out, err := cmd.CombinedOutput(); err != nil {
			return errors.New(string(out))
		}
	}
	return nil
}

func checkGitVersion(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "git", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.New(string(out))
	}
	// looks like "git version 2.38.1"
	r := regexp.MustCompile("git version ([0-9.]+)")
	matches := r.FindStringSubmatch(string(out))
	if len(matches) > 0 {
		comparison := semver.Compare("2.28", matches[1])
		return comparison >= 1, nil
	}
	return true, nil
}

// GitChecks prints warnings about uncommitted tracked and untracked files.
func GitChecks(ctx context.Context, l log.Logger, appPath string) error {
	l.Info(ctx, "Checking for uncommitted changes...")
	cmd := exec.Command("git", "status", "--porcelain=v2")
	cmd.Dir = appPath
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
		return fmt.Errorf("unable to proceed with deployment because of uncommitted changes")
	}
	l.Infof(ctx, "\t%s No uncommitted changes!", l.SuccessfulCheck())
	return nil
}

func GitMainBranch(branch string) bool {
	switch branch {
	case "main", "master":
		return true
	}

	return false
}

// GetGitSha will return the latest gitSha that will be used to create an application.
func GetGitSha(appPath string) (string, error) {
	// Gets latest git sha
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = appPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func GetGitBranch(appPath string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = appPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func GetTurbineResponseFromOutput(output string) (string, error) {
	r := regexp.MustCompile("turbine-response: ([^\n]*)")
	match := r.FindStringSubmatch(output)
	if match == nil || len(match) < 2 {
		return "", fmt.Errorf("output is formatted unexpectedly")
	}
	return match[1], nil
}

func RunCMD(ctx context.Context, logger log.Logger, cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		logger.Errorf(ctx, err.Error())
		return err
	}

	if err := cmd.Wait(); err != nil {
		logger.Errorf(ctx, err.Error())
		return err
	}
	return nil
}

// RunCmdWithErrorDetection checks exit codes and stderr for failures and logs on success.
func RunCmdWithErrorDetection(ctx context.Context, cmd *exec.Cmd, l log.Logger) (string, error) {
	stdout, stderr := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	stdOutMsg := stdout.String()
	stdErrMsg := stderr.String()

	if err != nil || stdErrMsg != "" {
		var errMsg, errLog string
		if err != nil {
			errMsg = err.Error()
		}
		if stdErrMsg != "" {
			// ignore most npm messages
			errorLogs := trimNonNpmErrorLines(stdErrMsg)
			if len(errorLogs) > 0 {
				errLog = errorLogs
			}
		} else if errMsg != "" {
			errLog = errMsg
		}
		if errLog != "" {
			if stdOutMsg != "" {
				l.Info(ctx, "\n"+stdOutMsg+"\n")
			}
			return "", errors.New(errLog)
		}
	}
	return stdOutMsg, nil
}

func trimNonNpmErrorLines(output string) string {
	ignoreThese := []string{"npm info", "npm timing", "npm http", "npm notice", "npm warn"}
	allLines := strings.Split(output, "\n")
	errorLines := []string{}

	for _, line := range allLines {
		skip := false
		for _, ignore := range ignoreThese {
			if strings.HasPrefix(line, ignore) {
				skip = true
				break
			}
		}
		if !skip {
			errorLines = append(errorLines, line)
		}
	}

	return strings.Join(errorLines, "\n")
}

// SwitchToAppDirectory switches temporarily to the application's directory.
func SwitchToAppDirectory(appPath string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return pwd, err
	}
	return pwd, os.Chdir(appPath)
}

func CleanupDockerfile(logger log.Logger, path string) {
	logger.StartSpinner("\t", fmt.Sprintf("Removing Dockerfile created for your application in %s...", path))
	// We clean up Dockerfile as last step

	if err := os.Remove(filepath.Join(path, "Dockerfile")); err != nil {
		logger.StopSpinnerWithStatus(fmt.Sprintf("Unable to remove Dockerfile: %v", err), log.Failed)
	} else {
		logger.StopSpinnerWithStatus("Dockerfile removed", log.Successful)
	}
}
