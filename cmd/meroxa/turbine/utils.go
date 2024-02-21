package turbine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/v2/pkg/ir"
)

const (
	GoLang     = "golang"
	JavaScript = "javascript"
	NodeJs     = "nodejs"
	Python     = "python"
	Python3    = "python3"
	Ruby       = "ruby"

	grpcFuncCollectionErr = "invalid ProcessCollectionRequest.Collection: embedded message failed validation | " +
		"caused by: invalid Collection.Name: value length must be at least 1 runes"
	grpcDestCollectionErr = "invalid WriteCollectionRequest.SourceCollection: embedded message failed validation | " +
		"caused by: invalid Collection.Name: value length must be at least 1 runes"
	missingSourceCollectionErr = `missing source or source collection, please ensure that you have configured your source correctly:
https://docs.meroxa.com/turbine/troubleshooting#troubleshooting-checklist"`
)

type AppConfig struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Language    ir.Lang           `json:"language"`
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

// GetLangFromAppJSON returns specified language in users' app.json.
func GetLangFromAppJSON(l log.Logger, pwd string) (ir.Lang, error) {
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
func GetAppNameFromAppJSON(l log.Logger, pwd string) (string, error) {
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
	return WriteConfigFile(pwd, appConfig)
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
	return WriteConfigFile(pwd, appConfig)
}

// ReadConfigFile will read the content of an app.json based on path.
func ReadConfigFile(appPath string) (*AppConfig, error) {
	var appConfig AppConfig

	if prefetched == nil || os.Getenv("UNIT_TEST") != "" {
		appConfigPath := path.Join(appPath, "app.json")
		appConfigBytes, err := os.ReadFile(appConfigPath)
		if err != nil {
			return &appConfig, fmt.Errorf("could not find an app.json file on path %q."+
				" Try a different value for `--path`", appPath)
		}
		if err := json.Unmarshal(appConfigBytes, &appConfig); err != nil {
			return &appConfig, err
		}
		prefetched = &appConfig
	}

	return prefetched, nil
}

func WriteConfigFile(appPath string, cfg *AppConfig) error {
	appConfigPath := path.Join(appPath, "app.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(appConfigPath, data, 0o664)
	if err != nil {
		return fmt.Errorf("%v\n"+
			"unable to update app.json file on path %q. Maybe try using a different value for `--path`", err, appPath)
	}
	return nil
}

func GetTurbineResponseFromOutput(output string) (string, error) {
	r := regexp.MustCompile("turbine-response: ([^\r\n]*)")
	match := r.FindStringSubmatch(output)
	if match == nil || len(match) < 2 {
		return "", fmt.Errorf("output is formatted unexpectedly")
	}

	trimmed := strings.TrimSpace(match[1])
	return trimmed, nil
}

// RunCmdWithErrorDetection checks exit codes and stderr for failures and logs on success.
func RunCmdWithErrorDetection(ctx context.Context, cmd *exec.Cmd, logger log.Logger) (string, error) {
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
				errLog = stdOutMsg + errLog
			}

			if strings.Contains(errLog, "rpc error") {
				logger.Debug(ctx, errLog)
				errLog = clarifyGrpcErrors(errLog)
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

// TODO: remove this and refactor error handling from turbine-core grpc requests
// This is needed temporarily to provide a more actionable error when the app has no sources defined.
// Longer-term, we need to have more specific, actionable errors rather than grpc validation
// errors surfaced to CLI output.
func clarifyGrpcErrors(errLog string) string {
	switch {
	case strings.Contains(errLog, grpcFuncCollectionErr):
		return missingSourceCollectionErr
	case strings.Contains(errLog, grpcDestCollectionErr):
		return missingSourceCollectionErr
	}
	return errLog
}

// SwitchToAppDirectory switches temporarily to the application's directory.
func SwitchToAppDirectory(appPath string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return pwd, err
	}
	return pwd, os.Chdir(appPath)
}

func UploadFile(ctx context.Context, logger log.Logger, filePath, url string) error {
	logger.StartSpinner("\t", "Uploading source...")

	var clientErr error
	var res *http.Response
	client := &http.Client{}

	retries := 3
	for retries > 0 {
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

		res, clientErr = client.Do(req)
		if clientErr != nil || (res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices) {
			retries--
		} else {
			break
		}
	}

	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	if clientErr != nil || (res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices) {
		logger.StopSpinnerWithStatus("\t Failed to upload build source file", log.Failed)
		if clientErr == nil {
			clientErr = fmt.Errorf("upload failed: %s", res.Status)
		}
		return clientErr
	}

	logger.StopSpinnerWithStatus("Source uploaded", log.Successful)
	return nil
}
