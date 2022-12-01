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
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
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
func (t *turbinePyCLI) Deploy(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error) {
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
	cmd.Env = append(
		cmd.Env,
		fmt.Sprintf("MEROXA_ACCESS_TOKEN=%s", accessToken),
		fmt.Sprintf("%s=%s", utils.AccountUUIDEnvVar, accountUUID))

	output, err = utils.RunCmdWithErrorDetection(ctx, cmd, t.logger)
	if err != nil {
		if strings.Contains(err.Error(), "unrecognized arguments") {
			return deploymentSpec, errors.New(utils.IncompatibleTurbineVersion)
		}

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
func (t *turbinePyCLI) GetResources(ctx context.Context, appName string) ([]utils.ApplicationResource, error) {
	var resources []utils.ApplicationResource

	cmd := exec.CommandContext(ctx, "turbine-py", "listResources", t.appPath)
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

func (t *turbinePyCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(t.appPath)
}

func (t *turbinePyCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}

func (t *turbinePyCLI) CreateDockerfile(ctx context.Context, appName string) (string, error) {
	cmd := exec.CommandContext(ctx, "turbine-py", "clibuild", t.appPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to create Dockerfile; %s", string(output))
	}
	r := regexp.MustCompile("^turbine-response: ([^\n]*)")
	match := r.FindStringSubmatch(string(output))
	if match == nil || len(match) < 2 {
		return "", fmt.Errorf("unable to create Dockerfiles; %s", string(output))
	}
	t.tmpPath = match[1]
	return t.tmpPath, err
}

func (t *turbinePyCLI) CleanUpBuild(ctx context.Context) {
	// cleanUpPythonTempBuildLocation removes any artifacts in the temporary directory.
	t.logger.StartSpinner("\t", fmt.Sprintf("Removing artifacts from building the Python Application at %s...", t.tmpPath))

	cmd := exec.CommandContext(ctx, "turbine-py", "cliclean", t.tmpPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.logger.StopSpinnerWithStatus(fmt.Sprintf("\t Failed to clean up artifacts at %s: %v %s", t.tmpPath, err, output), log.Failed)
	} else {
		t.logger.StopSpinnerWithStatus("Removed artifacts from building", log.Successful)
	}
}

func (t *turbinePyCLI) SetupForDeploy(ctx context.Context) (func(), error) {
	return func() {}, nil
}
