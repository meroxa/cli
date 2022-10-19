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
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
	turbineGo "github.com/meroxa/turbine-go/deploy"
)

const (
	GoLang     = "golang"
	JavaScript = "javascript"
	NodeJs     = "nodejs"
	Python     = "python"
	Python3    = "python3"

	turbineJSVersion  = "1.3.1"
	isTrue            = "true"
	AccountUUIDEnvVar = "MEROXA_ACCOUNT_UUID"
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

func GitInit(ctx context.Context, appPath string) error {
	if appPath == "" {
		return errors.New("path is required")
	}

	cmd := exec.Command("git", "config", "--global", "init.defaultBranch", "main")
	cmd.Path = appPath
	_ = cmd.Run()

	cmd = exec.Command("git", "init", appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(output))
	}
	return nil
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
func createTarAndZipFile(src string, buf io.Writer) error {
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

func UploadSource(ctx context.Context, logger log.Logger, language, appPath, appName, url string) error {
	var err error

	if language == GoLang || language == JavaScript {
		logger.StartSpinner("\t", fmt.Sprintf("Creating Dockerfile before uploading source in %s", appPath))

		if language == GoLang {
			err = turbineGo.CreateDockerfile(appName, appPath)
		}

		if language == JavaScript {
			err = createJavascriptDockerfile(ctx, appPath)
		}

		if err != nil {
			return err
		}
		defer func() {
			logger.StartSpinner("\t", fmt.Sprintf("Removing Dockerfile created for your application in %s...", appPath))
			// We clean up Dockerfile as last step
			err = os.Remove(filepath.Join(appPath, "Dockerfile"))
			if err != nil {
				logger.StopSpinnerWithStatus(fmt.Sprintf("Unable to remove Dockerfile: %v", err), log.Failed)
			} else {
				logger.StopSpinnerWithStatus("Dockerfile removed", log.Successful)
			}
		}()
		logger.StopSpinnerWithStatus("Dockerfile created", log.Successful)
	}

	dFile := fmt.Sprintf("turbine-%s.tar.gz", appName)

	var buf bytes.Buffer
	err = createTarAndZipFile(appPath, &buf)
	if err != nil {
		return err
	}

	logger.StartSpinner("\t", fmt.Sprintf(" Creating %q in %q to upload to our build service...", appPath, dFile))

	fileToWrite, err := os.OpenFile(dFile, os.O_CREATE|os.O_RDWR, os.FileMode(0777)) //nolint:gomnd
	defer func(fileToWrite *os.File) {
		err = fileToWrite.Close()
		if err != nil {
			panic(err.Error())
		}

		// remove .tar.gz file
		logger.StartSpinner("\t", fmt.Sprintf(" Removing %q...", dFile))
		removeErr := os.Remove(dFile)
		if removeErr != nil {
			logger.StopSpinnerWithStatus(fmt.Sprintf("\t Something went wrong trying to remove %q", dFile), log.Failed)
		} else {
			logger.StopSpinnerWithStatus(fmt.Sprintf("Removed %q", dFile), log.Successful)
		}

		if language == Python {
			cleanUpPythonTempBuildLocation(ctx, logger, appPath)
		}
	}(fileToWrite)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fileToWrite, &buf); err != nil {
		return err
	}
	logger.StopSpinnerWithStatus(fmt.Sprintf("%q successfully created in %q", dFile, appPath), log.Successful)

	return uploadFile(ctx, logger, dFile, url)
}

func uploadFile(ctx context.Context, logger log.Logger, filePath, url string) error {
	logger.StartSpinner("\t", " Uploading source...")

	fh, err := os.Open(filePath)
	if err != nil {
		logger.StopSpinnerWithStatus("\t Failed to open source file", log.Failed)
		return err
	}
	defer func(fh *os.File) {
		fh.Close()
	}(fh)

	req, err := http.NewRequestWithContext(ctx, "PUT", url, fh)
	if err != nil {
		logger.StopSpinnerWithStatus("\t Failed to make new request", log.Failed)
		return err
	}

	fi, err := fh.Stat()
	if err != nil {
		logger.StopSpinnerWithStatus("\t Failed to stat source file", log.Failed)
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	req.ContentLength = fi.Size()

	client := &http.Client{}

	retries := 3
	var res *http.Response
	for retries > 0 {
		res, err = client.Do(req)
		if err != nil {
			retries -= 1
		} else {
			break
		}
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		logger.StopSpinnerWithStatus("\t Failed to upload build source file", log.Failed)
		return err
	}

	logger.StopSpinnerWithStatus("Source uploaded", log.Successful)
	return nil
}

func RunTurbineJS(ctx context.Context, params ...string) (cmd *exec.Cmd) {
	args := getTurbineJSBinary(params)
	return executeTurbineJSCommand(ctx, args)
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

func executeTurbineJSCommand(ctx context.Context, params []string) *exec.Cmd {
	return exec.CommandContext(ctx, params[0], params[1:]...) //nolint:gosec
}

func createJavascriptDockerfile(ctx context.Context, appPath string) error {
	cmd := RunTurbineJS(ctx, "clibuild", appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to create Dockerfile at %s; %s", appPath, string(output))
	}

	return err
}

// cleanUpPythonTempBuildLocation removes any artifacts in the temporary directory.
func cleanUpPythonTempBuildLocation(ctx context.Context, logger log.Logger, appPath string) {
	logger.StartSpinner("\t", fmt.Sprintf(" Removing artifacts from building the Python Application at %s...", appPath))

	cmd := exec.CommandContext(ctx, "turbine-py", "cliclean", appPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		logger.StopSpinnerWithStatus(fmt.Sprintf("\t Failed to clean up artifacts at %s: %v %s", appPath, err, output), log.Failed)
	} else {
		logger.StopSpinnerWithStatus("Removed artifacts from building", log.Successful)
	}

	if err != nil {
		fmt.Printf("unable to clean up Meroxa Application at %s; %s", appPath, string(output))
	}
}

func GetTurbineResponseFromOutput(output string) (string, error) {
	r := regexp.MustCompile("turbine-response: ([^\n]*)")
	match := r.FindStringSubmatch(output)
	if match == nil || len(match) < 2 {
		return "", fmt.Errorf("output is formatted unexpectedly")
	}
	return match[1], nil
}
