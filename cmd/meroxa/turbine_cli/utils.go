package turbinecli

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

const turbineJSVersion = "1.0.0"

type AppConfig struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Language    string            `json:"language"`
	Resources   map[string]string `json:"resources"`
	Vendor      string            `json:"vendor"`
	ModuleInit  string            `json:"module_init"`
}

var prefetched *AppConfig
var isTrue = "true"

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

// GetLang will return language defined either by `--lang` or the one defined by user in the app.json.
func GetLang(ctx context.Context, l log.Logger, flag, pwd string) (string, error) {
	if flag != "" {
		return flag, nil
	}

	lang, err := GetLangFromAppJSON(ctx, l, pwd)
	if err != nil {
		return "", err
	} else if lang == "" {
		return "", fmt.Errorf("flag --lang is required unless lang is specified in your app.json")
	}
	return lang, nil
}

// GetLangFromAppJSON returns specified language in users' app.json.
func GetLangFromAppJSON(ctx context.Context, l log.Logger, pwd string) (string, error) {
	l.StartSpinner("\t", " Determining application language from app.json...")
	appConfig, err := readConfigFile(pwd)
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
	appConfig, err := readConfigFile(pwd)
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
	appConfig, err := readConfigFile(pwd)
	if err != nil {
		return err
	}
	appConfig.ModuleInit = isTrue
	if skipInit {
		appConfig.ModuleInit = "false"
	}
	err = writeConfigFile(pwd, appConfig) // will never be programmatically read again, but a marker of what turbine did
	return err
}

// SetVendorInAppJSON returns whether to vendor modules.
func SetVendorInAppJSON(pwd string, vendor bool) error {
	appConfig, err := readConfigFile(pwd)
	if err != nil {
		return err
	}
	appConfig.Vendor = "false"
	if vendor {
		appConfig.Vendor = isTrue
	}
	err = writeConfigFile(pwd, appConfig) // will never be programmatically read again, but a marker of what turbine did
	return err
}

// readConfigFile will read the content of an app.json based on path.
func readConfigFile(appPath string) (AppConfig, error) {
	var appConfig AppConfig

	if prefetched == nil {
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

func writeConfigFile(appPath string, cfg AppConfig) error {
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
	l.StartSpinner("", "Validating branch...")
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

func GoInit(ctx context.Context, l log.Logger, appPath string, skipInit, vendor bool) error {
	l.StartSpinner("\t", "Running golang module initializing...")
	skipLog := "skipping go module initialization\n\tFor guidance, visit " +
		"https://docs.meroxa.com/beta-overview#go-mod-init-for-a-new-golang-turbine-data-application"

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}
	if goPath == "" {
		l.StopSpinnerWithStatus("$GOPATH not set up; "+skipLog, log.Warning)
		return nil
	}
	i := strings.Index(appPath, goPath+"/src")
	if i == -1 || i != 0 {
		l.StopSpinnerWithStatus(fmt.Sprintf("%s is not under $GOPATH/src; %s", appPath, skipLog), log.Warning)
		return nil
	}

	// temporarily switching to the app's directory
	pwd, err := switchToAppDirectory(appPath)
	if err != nil {
		l.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	// initialize the user's module
	err = SetModuleInitInAppJSON(appPath, skipInit)
	if err != nil {
		l.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	err = modulesInit(l, appPath, skipInit, vendor)
	if err != nil {
		l.StopSpinnerWithStatus("\t", log.Failed)
		return err
	}

	return os.Chdir(pwd)
}

func modulesInit(l log.Logger, appPath string, skipInit, vendor bool) error {
	if skipInit {
		return nil
	}

	cmd := exec.Command("go", "mod", "init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		l.StopSpinnerWithStatus(fmt.Sprintf("\t%s", string(output)), log.Failed)
		return err
	}
	successLog := "go mod init succeeded"
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		successLog += fmt.Sprintf(" (while assuming GOPATH is %s)", build.Default.GOPATH)
	}
	l.StopSpinnerWithStatus(successLog+"!", log.Successful)
	l.StartSpinner("\t", "Getting latest turbine-go and turbine-go/running dependencies...")
	cmd = exec.Command("go", "get", "github.com/meroxa/turbine-go")
	output, err = cmd.CombinedOutput()
	if err != nil {
		l.StopSpinnerWithStatus(fmt.Sprintf("\t%s", string(output)), log.Failed)
		return err
	}
	cmd = exec.Command("go", "get", "github.com/meroxa/turbine-go/runner")
	output, err = cmd.CombinedOutput()
	if err != nil {
		l.StopSpinnerWithStatus(fmt.Sprintf("\t%s", string(output)), log.Failed)
		return err
	}
	l.StopSpinnerWithStatus("Downloaded latest turbine-go and turbine-go/running dependencies successfully!", log.Successful)

	// download dependencies
	err = SetVendorInAppJSON(appPath, vendor)
	if err != nil {
		return err
	}
	depsLog := "Downloading dependencies"
	cmd = exec.Command("go", "mod", "download")
	if vendor {
		depsLog += " to vendor"
		cmd = exec.Command("go", "mod", "vendor")
	}
	depsLog += "..."
	l.StartSpinner("\t", depsLog)
	output, err = cmd.CombinedOutput()
	if err != nil {
		l.StopSpinnerWithStatus(fmt.Sprintf("\t%s", string(output)), log.Failed)
		return err
	}
	l.StopSpinnerWithStatus("Downloaded all other dependencies successfully!", log.Successful)
	return nil
}

// switchToAppDirectory switches temporarily to the application's directory.
func switchToAppDirectory(appPath string) (string, error) {
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
	pwd, err := switchToAppDirectory(filepath.Dir(src))
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

func RunTurbineJS(params ...string) (cmd *exec.Cmd) {
	args := getTurbineJSBinary(params)
	return executeTurbineJSCommand(args)
}

func getTurbineJSBinary(params []string) []string {
	shouldUseLocalTurbineJS := global.GetLocalTurbineJSSetting()
	turbineJSBinary := fmt.Sprintf("@meroxa/turbine-js-cli@%s", turbineJSVersion)
	if shouldUseLocalTurbineJS == isTrue {
		turbineJSBinary = "turbine-js-cli"
	}
	args := []string{"npx", "--yes", turbineJSBinary}
	args = append(args, params...)
	return args
}

func executeTurbineJSCommand(params []string) *exec.Cmd {
	return exec.Command(params[0], params[1:]...) //nolint:gosec
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
