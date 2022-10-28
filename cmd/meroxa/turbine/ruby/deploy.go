package turbinerb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/meroxa/cli/cmd/meroxa/global"
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

// NeedsToBuild determines if the app has functions or not.
func (t *turbineRbCLI) NeedsToBuild(ctx context.Context, appName string) (bool, error) {
	//todo:
	// - run local server in "info" mode
	// - check if "process" is called
	return false, errors.New("unimplemented")
}

// Deploy creates Application entities.
func (t *turbineRbCLI) Deploy(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error) {
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
func (t *turbineRbCLI) GetResources(ctx context.Context, appName string) ([]utils.ApplicationResource, error) {
	return nil, errors.New("unimplemented")
}

func (t *turbineRbCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(t.appPath)
}

func (t *turbineRbCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}

func (t *turbineRbCLI) UploadSource(ctx context.Context, appName, tempPath, url string) error {
	return utils.UploadSource(ctx, t.logger, utils.Ruby, tempPath, appName, url)
}
