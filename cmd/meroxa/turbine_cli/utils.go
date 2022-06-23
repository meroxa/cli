package turbinecli

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

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

const turbineJSVersion = "0.2.1"

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
func GetLang(flag, pwd string) (string, error) {
	if flag != "" {
		return flag, nil
	}

	lang, err := GetLangFromAppJSON(pwd)
	if err != nil {
		return "", err
	} else if lang == "" {
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
	l.Info(ctx, "Checking for uncommitted changes...")
	// temporarily switching to the app's directory
	pwd, err := switchToAppDirectory(appPath)
	if err != nil {
		return err
	}

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
	l.Infof(ctx, "\t%s No uncommitted changes!", l.SuccessfulCheck())
	return os.Chdir(pwd)
}

// GetPipelineUUID parses the deploy output when it was successful to determine the pipeline UUID to create.
func GetPipelineUUID(output string) (string, error) {
	// Example output:
	// 2022/03/16 13:21:36 pipeline created: "turbine-pipeline-simple" ("049760a8-a3d2-44d9-b326-0614c09a3f3e").
	re := regexp.MustCompile(`pipeline:."turbine-pipeline-[a-z0-9-_]+".(\([^)]*\))`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 { //nolint:gomnd
		return "", fmt.Errorf("pipeline UUID not found")
	}
	return strings.Trim(matches[1], "()\""), nil
}

// ValidateBranch validates the deployment is being performed from one of the allowed branches.
func ValidateBranch(ctx context.Context, l log.Logger, appPath string) error {
	l.Info(ctx, "Validating branch...")
	// temporarily switching to the app's directory
	pwd, err := switchToAppDirectory(appPath)
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
	l.Infof(ctx, "\t%s Deployment allowed from %s branch!", l.SuccessfulCheck(), branchName)
	err = os.Chdir(pwd)
	if err != nil {
		return err
	}
	return nil
}

// GetGitSha will return the latest gitSha that will be used to create an application.
func GetGitSha(appPath string) (string, error) {
	// temporarily switching to the app's directory
	pwd, err := switchToAppDirectory(appPath)
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

func GoInit(ctx context.Context, l log.Logger, appPath string, skipInit, vendor bool) error {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		l.Warnf(ctx, "$GOPATH not set up; skipping go module initialization")
		return nil
	}
	i := strings.Index(appPath, goPath)
	if i == -1 || i != 0 {
		l.Warnf(ctx, "%s is not under $GOPATH; skipping go module initialization", appPath)
		return nil
	}

	// temporarily switching to the app's directory
	pwd, err := switchToAppDirectory(appPath)
	if err != nil {
		return err
	}

	// initialize the user's module
	err = SetModuleInitInAppJSON(appPath, skipInit)
	if err != nil {
		return err
	}
	if !skipInit {
		l.Info(ctx, "Initializing the application's go module...")
		cmd := exec.Command("go", "mod", "init")
		output, err := cmd.CombinedOutput()
		if err != nil {
			l.Error(ctx, string(output))
			return err
		}
		cmd = exec.Command("go", "get", "github.com/meroxa/turbine-go")
		output, err = cmd.CombinedOutput()
		if err != nil {
			l.Error(ctx, string(output))
			return err
		}
		cmd = exec.Command("go", "get", "github.com/meroxa/turbine-go/runner")
		output, err = cmd.CombinedOutput()
		if err != nil {
			l.Error(ctx, string(output))
			return err
		}

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
		l.Info(ctx, depsLog)
		output, err = cmd.CombinedOutput()
		if err != nil {
			l.Error(ctx, string(output))
			return err
		}
	}

	return os.Chdir(pwd)
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
		errLog = stdOutMsg
		if stdErrMsg != "" {
			errLog += stdErrMsg
		} else if errMsg != "" {
			errLog = errMsg
		}
		return "", errors.New(errLog)
	}
	if stdOutMsg != "" {
		l.Info(ctx, stdOutMsg)
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
		if fi.IsDir() && ((fi.Name() == ".git") || (fi.Name() == "fixtures")) {
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
	turbineJSBinary := fmt.Sprintf("@meroxa/turbine-js@%s", turbineJSVersion)
	if shouldUseLocalTurbineJS == isTrue {
		turbineJSBinary = "turbine"
	}
	args := []string{"npx", "--yes", turbineJSBinary}
	args = append(args, params...)
	return args
}

func executeTurbineJSCommand(params []string) *exec.Cmd {
	return exec.Command(params[0], params[1:]...) // nolint:gosec
}
