package turbineGo

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
)

type Deploy struct {
	DockerHubUserNameEnv    string
	DockerHubAccessTokenEnv string
	LocalDeployment         bool
}

// runDeployApp runs the binary previously built with the `--deploy` flag which should create all necessary resources
func runDeployApp(ctx context.Context, l log.Logger, appPath, appName, imageName string) error {
	l.Infof(ctx, "Deploying application %q...", appName)
	var cmd *exec.Cmd

	if imageName != "" {
		cmd = exec.Command(appPath+"/"+appName, "--deploy", "--imagename", imageName) // nolint:gosec
	} else {
		cmd = exec.Command(appPath+"/"+appName, "--deploy")
	}

	accessToken, refreshToken, err := global.GetUserToken()
	if err != nil {
		return err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("ACCESS_TOKEN=%s", accessToken), fmt.Sprintf("REFRESH_TOKEN=%s", refreshToken))

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Errorf(ctx, string(stdout))
		return err
	}

	l.Info(ctx, string(stdout))
	l.Info(ctx, "deploy complete!")
	return nil
}

// needsToBuild reads from the Turbine application to determine whether it needs to be built or not
// this is currently based on the number of functions
func needsToBuild(appPath, appName string) (bool, error) {
	cmd := exec.Command(appPath+"/"+appName, "--listfunctions")

	accessToken, refreshToken, err := global.GetUserToken()
	if err != nil {
		return false, err
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("ACCESS_TOKEN=%s", accessToken), fmt.Sprintf("REFRESH_TOKEN=%s", refreshToken))

	re := regexp.MustCompile(`\[(.*?)]`)
	stdout, err := cmd.CombinedOutput()

	// stdout is expected as `"2022/03/14 17:33:06 available functions: []` where within [], there will be each function
	hasFunctions := len(re.FindAllString(string(stdout), -1)) > 1
	return hasFunctions, nil
}

// Deploy takes care of all the necessary steps to deploy a Turbine application
//	1. Build binary
//	2. Build image
//	3. Push image
//	4. Run Turbine deploy
func (gd *Deploy) Deploy(ctx context.Context, appPath string, l log.Logger) error {
	var fqImageName string
	appName := path.Base(appPath)

	err := buildBinary(ctx, l, appPath, appName, true)
	if err != nil {
		return err
	}

	// check for image instances
	if ok, err := needsToBuild(appPath, appName); ok {
		if err != nil {
			l.Errorf(ctx, err.Error())
			return err
		}

		if gd.LocalDeployment {
			fqImageName, err = gd.getDockerImageName(ctx, l, appPath, appName)
			if err != nil {
				return err
			}
		} else {
			fqImageName, err = gd.getPlatformImage(ctx, l)
		}
	}

	// creates all resources
	err = runDeployApp(ctx, l, appPath, appName, fqImageName)
	if err != nil {
		l.Errorf(ctx, "unable to deploy app; %s", err)
		return err
	}
	return nil
}
