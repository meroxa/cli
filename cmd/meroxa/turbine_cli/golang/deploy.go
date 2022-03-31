package turbinego

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	turbinecli "github.com/meroxa/cli/cmd/meroxa/turbine_cli"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

type Deploy struct {
	DockerHubUserNameEnv    string
	DockerHubAccessTokenEnv string
	LocalDeployment         bool
}

// RunDeployApp runs the binary previously built with the `--deploy` flag which should create all necessary resources.
func RunDeployApp(ctx context.Context, l log.Logger, appPath, appName, imageName string) (string, error) {
	l.Infof(ctx, "Deploying application %q...", appName)
	var cmd *exec.Cmd

	if imageName != "" {
		cmd = exec.Command(appPath+"/"+appName, "--deploy", "--imagename", imageName) // nolint:gosec
	} else {
		cmd = exec.Command(appPath+"/"+appName, "--deploy") // nolint:gosec
	}

	accessToken, refreshToken, err := global.GetUserToken()
	if err != nil {
		return "", err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("ACCESS_TOKEN=%s", accessToken), fmt.Sprintf("REFRESH_TOKEN=%s", refreshToken))

	stdout, err := turbinecli.RunCmdWithErrorDetection(ctx, cmd, l)
	if err == nil {
		l.Info(ctx, "deploy complete!")
	}
	return stdout, err
}

// NeedsToBuild reads from the Turbine application to determine whether it needs to be built or not
// this is currently based on the number of functions.
func NeedsToBuild(appPath, appName string) (bool, error) {
	cmd := exec.Command(appPath+"/"+appName, "--listfunctions") // nolint:gosec

	accessToken, refreshToken, err := global.GetUserToken()
	if err != nil {
		return false, err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("ACCESS_TOKEN=%s", accessToken), fmt.Sprintf("REFRESH_TOKEN=%s", refreshToken))

	// TODO: Implement in Turbine something that returns a boolean rather than having to regex its output
	re := regexp.MustCompile(`\[(.+?)]`)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	// stdout is expected as `"2022/03/14 17:33:06 available functions: []` where within [], there will be each function.
	hasFunctions := len(re.FindAllString(string(stdout), -1)) > 0

	return hasFunctions, nil
}
