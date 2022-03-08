package turbinecli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/global"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	turbine "github.com/meroxa/turbine/deploy"

	"github.com/meroxa/cli/log"
)

type GoDeploy struct {
	DockerHubUserNameEnv    string
	DockerHubAccessTokenEnv string
}

// buildGoApp will create a go binary with a specific name on a specific path.
func buildGoApp(ctx context.Context, l log.Logger, appPath, appName string, platform bool) error {
	var cmd *exec.Cmd

	if platform {
		cmd = exec.Command("go", "build", "--tags", "platform", "-o", appName, "./...")
	} else {
		cmd = exec.Command("go", "build", "-o", appName, "./...")
	}
	cmd.Dir = appPath

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, string(stdout))
		return err
	}

	return nil
}

// deployApp deploys the image previously built.
func deployApp(ctx context.Context, l log.Logger, appPath, appName, imageName string) error {
	l.Infof(ctx, "Deploying application %q...", appName)

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

func (gd *GoDeploy) getAuthConfig() string {
	dhUsername := gd.DockerHubUserNameEnv
	dhAccessToken := gd.DockerHubAccessTokenEnv
	authConfig := types.AuthConfig{
		Username:      dhUsername,
		Password:      dhAccessToken,
		ServerAddress: "https://index.docker.io/v1/",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	return base64.URLEncoding.EncodeToString(authConfigBytes)
}

func (gd *GoDeploy) pushImage(l log.Logger, imageName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	authConfig := gd.getAuthConfig()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120) // nolint:gomnd
	defer cancel()

	l.Infof(ctx, "pushing image %q to container registry", imageName)
	opts := types.ImagePushOptions{RegistryAuth: authConfig}
	rd, err := cli.ImagePush(ctx, imageName, opts)
	if err != nil {
		return err
	}
	defer func(rd io.ReadCloser) {
		err = rd.Close()
		if err != nil {
			l.Error(ctx, err.Error())
		}
	}(rd)

	_, err = io.Copy(os.Stdout, rd)
	if err != nil {
		return err
	}
	l.Info(ctx, "image successfully pushed to container registry!")

	return nil
}

func (*GoDeploy) buildImage(ctx context.Context, l log.Logger, pwd, imageName string) error {
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

// RunGoApp will build a go binary and will run it.
func RunGoApp(ctx context.Context, appPath string, l log.Logger) error {
	// grab current location to use it as project name
	appName := path.Base(appPath)

	// building is a requirement prior to running for go apps
	err := buildGoApp(ctx, l, appPath, appName, false)
	if err != nil {
		return err
	}

	err = os.Chdir(appPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("./" + appName) //nolint:gosec

	out, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, err.Error())
	}

	l.Info(ctx, string(out))
	return nil
}

func (gd *GoDeploy) DeployGoApp(ctx context.Context, appPath string, l log.Logger) error {
	appName := path.Base(appPath)
	fqImageName := gd.DockerHubUserNameEnv + "/" + appName
	err := gd.buildImage(ctx, l, appPath, fqImageName)
	if err != nil {
		l.Errorf(ctx, "unable to build image; %q\n%s", fqImageName, err)
	}

	err = gd.pushImage(l, fqImageName)
	if err != nil {
		l.Errorf(ctx, "unable to push image; %q\n%s", fqImageName, err)
	}

	err = buildGoApp(ctx, l, appPath, appName, true)
	if err != nil {
		return err
	}

	err = deployApp(ctx, l, appPath, appName, fqImageName)
	if err != nil {
		l.Errorf(ctx, "unable to deploy app; %s", err)
	}

	return nil
}
