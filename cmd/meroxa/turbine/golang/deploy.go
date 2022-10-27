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
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/global"
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// Deploy runs the binary previously built with the `--deploy` flag which should create all necessary resources.
func (t *turbineGoCLI) Deploy(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error) {
	deploymentSpec := ""
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
	cmd := exec.CommandContext(ctx, path.Join(t.appPath, appName), args...) //nolint:gosec

	accessToken, refreshToken, err := global.GetUserToken()
	if err != nil {
		return deploymentSpec, err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(
		cmd.Env,
		fmt.Sprintf("ACCESS_TOKEN=%s", accessToken),
		fmt.Sprintf("REFRESH_TOKEN=%s", refreshToken),
		fmt.Sprintf("MEROXA_DEBUG=%s", os.Getenv("MEROXA_DEBUG")),
		fmt.Sprintf("%s=%s", utils.AccountUUIDEnvVar, accountUUID))

	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "flag provided but not defined") {
			return deploymentSpec, errors.New(utils.IncompatibleTurbineVersion)
		}
		return deploymentSpec, errors.New(string(output))
	}

	if specVersion != "" {
		deploymentSpec, err = utils.GetTurbineResponseFromOutput(string(output))
		if err != nil {
			err = fmt.Errorf(
				"unable to receive the deployment spec for the Meroxa Application at %s", t.appPath)
		}
	}
	return deploymentSpec, err
}

func (t *turbineGoCLI) GetResources(ctx context.Context, appName string) ([]utils.ApplicationResource, error) {
	var resources []utils.ApplicationResource

	cmd := exec.CommandContext(ctx, path.Join(t.appPath, appName), "--listresources") //nolint:gosec
	output, err := cmd.CombinedOutput()
	if err != nil {
		return resources, errors.New(string(output))
	}
	list, err := utils.GetTurbineResponseFromOutput(string(output))
	if err == nil && list != "" {
		// ignores any lines that are not intended to be part of the response
		output = []byte(list)
	}

	if err := json.Unmarshal(output, &resources); err != nil {
		// fall back if not json
		return utils.GetResourceNamesFromString(string(output)), nil
	}

	return resources, nil
}

// NeedsToBuild reads from the Turbine application to determine whether it needs to be built or not
// this is currently based on the number of functions.
func (t *turbineGoCLI) NeedsToBuild(ctx context.Context, appName string) (bool, error) {
	cmd := exec.CommandContext(ctx, path.Join(t.appPath, appName), "--listfunctions") //nolint:gosec

	accessToken, refreshToken, err := global.GetUserToken()
	if err != nil {
		return false, err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("ACCESS_TOKEN=%s", accessToken), fmt.Sprintf("REFRESH_TOKEN=%s", refreshToken))

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return false, fmt.Errorf("build failed")
	}
	list, err := utils.GetTurbineResponseFromOutput(string(output))
	if err == nil && list != "" {
		// ignores any lines that are not intended to be part of the response
		list = strings.Trim(list, "[]")
		if list == "" {
			return false, nil
		}

		hasFunctions := len(strings.Split(list, ",")) > 0
		return hasFunctions, nil
	}

	re := regexp.MustCompile(`\[(.+?)]`)
	// stdout is expected as `"2022/03/14 17:33:06 available functions: []` where within [], there will be each function.
	hasFunctions := len(re.FindAllString(string(output), -1)) > 0

	return hasFunctions, nil
}

func (t *turbineGoCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(t.appPath)
}

func (t *turbineGoCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}

func (t *turbineGoCLI) UploadSource(ctx context.Context, appName, tempPath, url string) error {
	return utils.UploadSource(ctx, t.logger, utils.GoLang, t.appPath, appName, url)
}
