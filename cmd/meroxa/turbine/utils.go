package turbine

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/meroxa/cli/log"
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

// getResourceNamesFromString provides backward compatibility with turbine-go
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
	l.StartSpinner("\t", " Determining application language from app.json...")
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
	l.StartSpinner("\t", " Reading application name from app.json...")
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

// GitChecks prints warnings about uncommitted tracked and untracked files.
func GitChecks(ctx context.Context, l log.Logger, appPath string) error {
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

// ValidateBranch validates the deployment is being performed from one of the allowed branches.
func ValidateBranch(ctx context.Context, l log.Logger, appPath string) error {
	l.StartSpinner("", " Validating branch...")
	branchName, err := GetGitBranch(appPath)
	if err != nil {
		return err
	}

	if GitMainBranch(branchName) {
		l.StopSpinnerWithStatus(fmt.Sprintf("Deployment allowed from %q branch", branchName), log.Successful)
		return nil
	}

	l.StopSpinnerWithStatus(fmt.Sprintf("deployment allowed only from \"main\" or \"master\" branch, not %q", branchName), log.Failed)
	return fmt.Errorf(`deployment allowed only from "main" or "master" branch, not %q`, branchName)
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

// SwitchToAppDirectory switches temporarily to the application's directory.
func SwitchToAppDirectory(appPath string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return pwd, err
	}
	return pwd, os.Chdir(appPath)
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
			errLog = stdErrMsg
		} else if errMsg != "" {
			errLog = errMsg
		}
		if stdOutMsg != "" {
			l.Info(ctx, "\n"+stdOutMsg+"\n")
		}
		return "", errors.New(errLog)
	}
	return stdOutMsg, nil
}

// CreateTarAndZipFile creates a .tar.gz file from `src` on current directory.
func CreateTarAndZipFile(src string, buf io.Writer) error {
	// Grab the directory we care about (app's directory)
	appDir := filepath.Base(src)

	// Change to parent's app directory
	pwd, err := SwitchToAppDirectory(filepath.Dir(src))
	if err != nil {
		return err
	}

	zipWriter := gzip.NewWriter(buf)
	tarWriter := tar.NewWriter(zipWriter)

	err = filepath.Walk(appDir, func(file string, fi os.FileInfo, err error) error {
		if shouldSkipDir(fi) {
			return filepath.SkipDir
		}
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(file)
		if err := tarWriter.WriteHeader(header); err != nil { //nolint:govet
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		if !fi.IsDir() {
			var data *os.File
			data, err = os.Open(file)
			defer func(data *os.File) {
				err = data.Close()
				if err != nil {
					panic(err.Error())
				}
			}(data)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := tarWriter.Close(); err != nil {
		return err
	}
	if err := zipWriter.Close(); err != nil {
		return err
	}

	return os.Chdir(pwd)
}

func shouldSkipDir(fi os.FileInfo) bool {
	if !fi.IsDir() {
		return false
	}

	switch fi.Name() {
	case ".git", "fixtures", "node_modules":
		return true
	}

	return false
}
