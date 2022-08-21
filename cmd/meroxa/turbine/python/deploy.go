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
	"github.com/meroxa/cli/cmd/meroxa/turbine"
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
func RunDeployApp(ctx context.Context, l log.Logger, path, imageName, appName, gitSha, specVersion string) error {
	args := []string{"clideploy", path, imageName, appName, gitSha}

	if specVersion != "" {
		args = append(args, specVersion)
	}

	cmd := exec.Command("turbine-py", args...)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	_, err = turbine.RunCmdWithErrorDetection(ctx, cmd, l)
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

// GetResources asks turbine for a list of resources used by the given app.
func GetResources(ctx context.Context, l log.Logger, appPath, appName string) ([]turbine.ApplicationResource, error) {
	var resources []turbine.ApplicationResource

	cmd := exec.Command("turbine-py", "listResources", appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return resources, errors.New(string(output))
	}

	if err := json.Unmarshal(output, &resources); err != nil {
		// fall back if not json
		return turbine.GetResourceNamesFromString(string(output)), nil
	}

	return resources, nil
}
