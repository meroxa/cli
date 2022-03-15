package turbineGo

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"

	turbine "github.com/meroxa/turbine/deploy"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
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

	// TODO: Check here if imageName is ""
	cmd := exec.Command(appPath+"/"+appName, "--deploy", "--imagename", imageName) // nolint:gosec
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

func (*Deploy) buildImage(ctx context.Context, l log.Logger, pwd, imageName string) error {
	l.Infof(ctx, "Building image %q located at %q", imageName, pwd)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		l.Errorf(ctx, "unable to init docker client; %s", err)
	}

	// Generate dockerfile
	err = turbine.CreateDockerfile(pwd)
	if err != nil {
		return err
	}

	// Read local Dockerfile
	tar, err := archive.TarWithOptions(pwd, &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{".git", "fixtures"},
	})
	if err != nil {
		l.Errorf(ctx, "unable to create tar; %s", err)
	}

	buildOptions := types.ImageBuildOptions{
		Context:    tar,
		Dockerfile: "Dockerfile",
		Remove:     true,
		Tags:       []string{imageName}}

	resp, err := cli.ImageBuild(
		ctx,
		tar,
		buildOptions,
	)
	if err != nil {
		l.Errorf(ctx, "unable to build docker image; %s", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			l.Errorf(ctx, "unable to close docker build response body; %s", err)
		}
	}(resp.Body)
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		l.Errorf(ctx, "unable to read image build response; %s", err)
	}
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

	// to determine whether we need to call POST /sources and POST /build
	if ok, err := needsToBuild(appPath, appName); ok {
		if err != nil {
			l.Errorf(ctx, err.Error())
			return err
		}

		//fqImageName will be eventually taken from the build endpoint
		fqImageName = gd.DockerHubUserNameEnv + "/" + appName

		err = gd.buildImage(ctx, l, appPath, fqImageName)
		if err != nil {
			l.Errorf(ctx, "unable to build image; %q\n%s", fqImageName, err)
			return err
		}

		// this will go away when using the new service
		err = gd.pushImage(l, fqImageName)
		if err != nil {
			l.Errorf(ctx, "unable to push image; %q\n%s", fqImageName, err)
			return err
		}

		// TODO: Keep polling for building statsu until it's successful
	}

	// creates all resources
	err = runDeployApp(ctx, l, appPath, appName, fqImageName)
	if err != nil {
		l.Errorf(ctx, "unable to deploy app; %s", err)
		return err
	}
	return nil
}
