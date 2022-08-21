package turbinego

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
)

// RunDeployApp runs the binary previously built with the `--deploy` flag which should create all necessary resources.
func RunDeployApp(ctx context.Context, l log.Logger, appPath, imageName, appName, gitSha string, specVersion string) error {
	args := []string{
		"--deploy",
		"--appname",
		appName,
		"--gitsha",
		gitSha,
	}

	if specVersion != "" {
		args = append(args, "--spec", specVersion)
	}

	if imageName != "" {
		args = append(args, "--imagename", imageName)
	}
	cmd := exec.Command(path.Join(appPath, appName), args...) //nolint:gosec

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

func GetResources(ctx context.Context, l log.Logger, appPath, appName string) ([]turbine.ApplicationResource, error) {
	var resources []turbine.ApplicationResource

	cmd := exec.Command(path.Join(appPath, appName), "--listresources") //nolint:gosec
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

// NeedsToBuild reads from the Turbine application to determine whether it needs to be built or not
// this is currently based on the number of functions.
func NeedsToBuild(appPath, appName string) (bool, error) {
	cmd := exec.Command(path.Join(appPath, appName), "--listfunctions") //nolint:gosec

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
