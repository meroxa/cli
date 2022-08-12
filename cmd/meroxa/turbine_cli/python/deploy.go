package turbinepy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/global"
	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"
	"github.com/meroxa/cli/log"
)

// BuildApp created the needed structure for a python app.
func BuildApp(path string) (string, error) {
	cmd := exec.Command("turbine-py", "clibuild", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to build Meroxa Application at %s; %s", path, string(output))
	}
	r := regexp.MustCompile("^turbine-response: ([^\n]*)")
	match := r.FindStringSubmatch(string(output))
	if match == nil || len(match) < 2 {
		return "", fmt.Errorf("unable to build Meroxa Application at %s; %s", path, string(output))
	}
	return match[1], err
}

// NeedsToBuild determines if the app has functions or not.
func NeedsToBuild(path string) (bool, error) {
	cmd := exec.Command("turbine-py", "hasFunctions", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err := fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			path,
			string(output))
		return false, err
	}
	r := regexp.MustCompile("^turbine-response: (true|false)\n")
	match := r.FindStringSubmatch(string(output))
	if match == nil || len(match) < 2 {
		err := fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			path,
			string(output))
		return false, err
	}
	return strconv.ParseBool(match[1])
}

// RunDeployApp creates Application entities.
func RunDeployApp(ctx context.Context, l log.Logger, path, imageName, appName, gitSha string) error {
	cmd := exec.Command("turbine-py", "clideploy", path, imageName, appName, gitSha)
  
	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	_, err = turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
	return err
}

// CleanUpApp removes any temporary artifacts in the temp directory.
func CleanUpApp(path string) (string, error) {
	cmd := exec.Command("turbine-py", "cliclean", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to clean up Meroxa Application at %s; %s", path, string(output))
	}
	return strings.TrimSpace(string(output)), err
}

// GetResourceNames asks turbine for a list of resources used by the given app.
func GetResourceNames(ctx context.Context, l log.Logger, appPath, appName string) ([]string, error) {
	var names []string

	cmd := exec.Command("turbine-py", "listResources", appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return names, errors.New(string(output))
	}

	var parsed []turbinecli.ApplicationResource
	if err := json.Unmarshal(output, &parsed); err != nil {
		// fall back if not json
		return getResourceNamesFromString(string(output)), nil
	}

	for i := range parsed {
		kv := parsed[i]
		if kv.Name != "" {
			names = append(names, kv.Name)
		}
	}
	return names, nil
}

// getResourceNamesFromString provides backward compatibility with turbine-go
// legacy resource listing format.
func getResourceNamesFromString(s string) []string {
	var names []string

	r := regexp.MustCompile(`\[(.+?)\]`)
	sliceString := r.FindStringSubmatch(s)
	if len(sliceString) > 0 {
		names = strings.Fields(sliceString[1])
	}

	return names
}
