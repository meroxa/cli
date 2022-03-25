package turbinego

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
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

// runDeployApp runs the binary previously built with the `--deploy` flag which should create all necessary resources.
func runDeployApp(ctx context.Context, l log.Logger, appPath, appName, imageName string) (string, error) {
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

// needsToBuild reads from the Turbine application to determine whether it needs to be built or not
// this is currently based on the number of functions.
func needsToBuild(appPath, appName string) (bool, error) {
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

// Deploy takes care of all the necessary steps to deploy a Turbine application
//	1. Build binary
//	2. Build image
//	3. Push image
//	4. Run Turbine deploy
func (gd *Deploy) Deploy(ctx context.Context, appPath string, l log.Logger) (string, error) {
	var fqImageName string
	appName := path.Base(appPath)

	err := buildBinary(ctx, l, appPath, appName, true)
	if err != nil {
		return "", err
	}

	// check for image instances
	var ok bool
	if ok, err = needsToBuild(appPath, appName); ok {
		if err != nil {
			l.Errorf(ctx, err.Error())
			return "", err
		}

		if gd.LocalDeployment {
			fqImageName, err = gd.getDockerImageName(ctx, l, appPath, appName)
			l.Infof(ctx, "fqImageName: %q", fqImageName)
			if err != nil {
				return "", err
			}
		} else {
			// fqImageName, err = gd.getPlatformImage(ctx, l).
			_, _ = gd.getPlatformImage(ctx, l)
			// Returns so it doesn't error the next step
			return "", fmt.Errorf("using build service not working at the moment, please use your own dockerhub credentials in the meantime")
		}
	}

	// creates all resources
	output, err := runDeployApp(ctx, l, appPath, appName, fqImageName)
	if err != nil {
		l.Errorf(ctx, "unable to deploy app; %s", err)
		return output, err
	}
	return output, nil
}
