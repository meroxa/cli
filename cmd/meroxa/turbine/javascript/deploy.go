package turbinejs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/meroxa/cli/cmd/meroxa/global"
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineJsCLI) NeedsToBuild(ctx context.Context, appName string) (bool, error) {
	cmd := utils.RunTurbineJS(ctx, "hasfunctions", t.appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			t.appPath,
			string(output))
		return false, err
	}

	isNeeded, err := utils.GetTurbineResponseFromOutput(string(output))
	if err != nil {
		err := fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			t.appPath,
			string(output))
		return false, err
	}
	return strconv.ParseBool(isNeeded)
}

func (t *turbineJsCLI) Deploy(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error) {
	var (
		output         string
		deploymentSpec string
	)
	params := []string{"clideploy", imageName, t.appPath, appName, gitSha}

	if specVersion != "" {
		params = append(params, specVersion)
	}

	cmd := utils.RunTurbineJS(ctx, params...)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return deploymentSpec, err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(
		cmd.Env,
		fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken),
		fmt.Sprintf("%s=%s", utils.AccountUUIDEnvVar, accountUUID))

	output, err = utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	if err != nil {
		return deploymentSpec, err
	}
	if specVersion != "" {
		deploymentSpec, err = utils.GetTurbineResponseFromOutput(output)
		if err != nil {
			err = fmt.Errorf(
				"unable to receive the deployment spec for the Meroxa Application at %s", t.appPath)
		}
	}
	return deploymentSpec, err
}

// GetResources asks turbine for a list of resources used by the given app.
func (t *turbineJsCLI) GetResources(ctx context.Context, appName string) ([]utils.ApplicationResource, error) {
	var resources []utils.ApplicationResource

	cmd := utils.RunTurbineJS(ctx, "listresources", t.appPath)
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

func (t *turbineJsCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(t.appPath)
}

func (t *turbineJsCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}

func (t *turbineJsCLI) UploadSource(ctx context.Context, appName, appPath, url string) error {
	return utils.UploadSource(ctx, t.logger, utils.JavaScript, t.appPath, appName, url)
}
