package turbinego

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

// RunDeployApp runs the binary previously built with the `--deploy` flag which should create all necessary resources.
func RunDeployApp(ctx context.Context, l log.Logger, appPath, imageName, appName, gitSha, specVersion string) error {
	var cmd *exec.Cmd

	args := []string{"--deploy", "--gitsha", gitSha}

	if specVersion != "" {
		args = append(args, "--spec", specVersion)
	}

	if imageName != "" {
		args = append(args, "--imagename", imageName)
	}

	cmd = exec.Command(appPath+"/"+appName, args...) //nolint:gosec

	accessToken, refreshToken, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("ACCESS_TOKEN=%s", accessToken), fmt.Sprintf("REFRESH_TOKEN=%s", refreshToken))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

// GetResourceNames asks turbine for a list of resources used by the given app.
func GetResourceNames(ctx context.Context, l log.Logger, appPath, appName string) ([]string, error) {
	var names []string

	cmd := exec.Command(appPath+"/"+appName, "--listresources") //nolint:gosec

	output, err := cmd.CombinedOutput()
	if err != nil {
		return names, errors.New(string(output))
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal(output, &parsed); err != nil {
		// fall back if not json
		return getResourceNamesFromString(string(output)), nil
	}

	for i := range parsed {
		kv := parsed[i]
		if _, ok := kv["Name"]; ok {
			names = append(names, kv["Name"].(string))
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

// NeedsToBuild reads from the Turbine application to determine whether it needs to be built or not
// this is currently based on the number of functions.
func NeedsToBuild(appPath, appName string) (bool, error) {
	cmd := exec.Command(appPath+"/"+appName, "--listfunctions") //nolint:gosec

	accessToken, refreshToken, err := global.GetUserToken()
	if err != nil {
		return false, err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("ACCESS_TOKEN=%s", accessToken), fmt.Sprintf("REFRESH_TOKEN=%s", refreshToken))

	// TODO: Implement in Turbine something that returns a boolean rather than having to regex its output
	re := regexp.MustCompile(`\[(.+?)]`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return false, fmt.Errorf("build failed")
	}

	// stdout is expected as `"2022/03/14 17:33:06 available functions: []` where within [], there will be each function.
	hasFunctions := len(re.FindAllString(string(output), -1)) > 0

	return hasFunctions, nil
}
