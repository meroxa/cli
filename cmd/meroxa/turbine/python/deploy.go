package turbinepy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/meroxa/cli/cmd/meroxa/global"
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// NeedsToBuild determines if the app has functions or not.
func (t *turbinePyCLI) NeedsToBuild(ctx context.Context, appName string) (bool, error) {
	cmd := exec.Command("turbine-py", "hasFunctions", t.appPath)
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

// Deploy creates Application entities.
func (t *turbinePyCLI) Deploy(ctx context.Context, imageName, appName, gitSha, specVersion string) (string, error) {
	var (
		output         string
		deploymentSpec string
	)
	args := []string{"clideploy", t.appPath, imageName, appName, gitSha}

	if specVersion != "" {
		args = append(args, specVersion)
	}

	cmd := exec.CommandContext(ctx, "turbine-py", args...)

	accessToken, _, err := global.GetUserToken()
	if err != nil {
		return deploymentSpec, err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken))

	output, err = utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	if err != nil {
		return deploymentSpec, err
	}
	if specVersion != "" {
		deploymentSpec, err = utils.GetTurbineResponseFromOutput(output)
		if err != nil {
			err = fmt.Errorf(
				"unable to receive the deployment spec for the Meroxa Application at %s has a Process", t.appPath)
		}
	}
	return deploymentSpec, err
}

// GetResources asks turbine for a list of resources used by the given app.
func (t *turbinePyCLI) GetResources(ctx context.Context, appName string) ([]utils.ApplicationResource, error) {
	var resources []utils.ApplicationResource

	cmd := exec.CommandContext(ctx, "turbine-py", "listResources", t.appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return resources, errors.New(string(output))
	}

	if err := json.Unmarshal(output, &resources); err != nil {
		// fall back if not json
		return utils.GetResourceNamesFromString(string(output)), nil
	}

	return resources, nil
}

func (t *turbinePyCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(t.appPath)
}

func (t *turbinePyCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}

func (t *turbinePyCLI) UploadSource(ctx context.Context, appName, tempPath, url string) error {
	return utils.UploadSource(ctx, t.logger, utils.Python, tempPath, appName, url)
}
