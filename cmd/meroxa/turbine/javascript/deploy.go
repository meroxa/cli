package turbinejs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/meroxa/cli/cmd/meroxa/global"
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

var TurbineJSVersion = "1.3.8"

func (t *turbineJsCLI) NeedsToBuild(ctx context.Context) (bool, error) {
	output, err := RunTurbineJS(ctx, t.logger, "hasfunctions", t.appPath)
	if err != nil {
		err = fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			t.appPath,
			err)
		return false, err
	}

	isNeeded, err := utils.GetTurbineResponseFromOutput(output)
	if err != nil {
		fmtErr := fmt.Errorf(
			"unable to determine if the Meroxa Application at %s has a Process; %s",
			t.appPath,
			err)
		return false, fmtErr
	}
	return strconv.ParseBool(isNeeded)
}

func (t *turbineJsCLI) GetDeploymentSpec(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error) {
	var (
		output         string
		deploymentSpec string
	)
	params := []string{"clideploy", imageName, t.appPath, appName, gitSha}

	if specVersion != "" {
		params = append(params, specVersion)
	}

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return deploymentSpec, err
	}
	output, err = RunTurbineJSWithEnvVars(
		ctx,
		t.logger,
		map[string]string{
			"MEROXA_ACCESS_TOKEN":   accessToken,
			utils.AccountUUIDEnvVar: accountUUID,
		},
		params...)
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
func (t *turbineJsCLI) GetResources(ctx context.Context) ([]utils.ApplicationResource, error) {
	var resources []utils.ApplicationResource

	output, err := RunTurbineJS(ctx, t.logger, "listresources", t.appPath)
	if err != nil {
		return resources, err
	}
	list, err := utils.GetTurbineResponseFromOutput(output)
	if err == nil && list != "" {
		// ignores any lines that are not intended to be part of the response
		output = list
	}

	if err := json.Unmarshal([]byte(output), &resources); err != nil {
		// fall back if not json
		return utils.GetResourceNamesFromString(output), nil
	}

	return resources, nil
}

func (t *turbineJsCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(ctx, t.appPath)
}

func (t *turbineJsCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}

func (t *turbineJsCLI) CreateDockerfile(ctx context.Context, _ string) (string, error) {
	_, err := RunTurbineJS(ctx, t.logger, "clibuild", t.appPath)
	if err != nil {
		return "", fmt.Errorf("unable to create Dockerfile at %s; %s", t.appPath, err)
	}

	return t.appPath, err
}

func (t *turbineJsCLI) CleanUpBuild(_ context.Context) {
	utils.CleanupDockerfile(t.logger, t.appPath)
}

func (t *turbineJsCLI) SetupForDeploy(_ context.Context, _ string) (func(), error) {
	return func() {}, nil
}
