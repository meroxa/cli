package apps

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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
	turbine "github.com/meroxa/turbine/deploy"
)

const (
	dockerHubUserNameEnv    = "DOCKER_HUB_USERNAME"
	dockerHubAccessTokenEnv = "DOCKER_HUB_ACCESS_TOKEN" // nolint:gosec
	goLang                  = "golang"
	javaScript              = "javascript"
	nodeJS                  = "nodejs"
)

func buildGoApp(ctx context.Context, l log.Logger, appPath, appName string, platform bool) error {
	var cmd *exec.Cmd

	if appName != "" {
		appName = appPath
	}

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

func deployApp(ctx context.Context, l log.Logger, appPath, appName, imageName string) error {
	l.Info(ctx, "deploying app...\n")

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

func buildImage(ctx context.Context, l log.Logger, pwd, imageName string) error {
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

func pushImage(l log.Logger, imageName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	authConfig := getAuthConfig()
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

func getAuthConfig() string {
	dhUsername := os.Getenv(dockerHubUserNameEnv)
	dhAccessToken := os.Getenv(dockerHubAccessTokenEnv)
	authConfig := types.AuthConfig{
		Username:      dhUsername,
		Password:      dhAccessToken,
		ServerAddress: "https://index.docker.io/v1/",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	return base64.URLEncoding.EncodeToString(authConfigBytes)
}

func prependAccount(imageName string) string {
	account := os.Getenv(dockerHubUserNameEnv)
	return account + "/" + imageName
}

func readConfigFile(appPath string) (AppConfig, error) {
	var appConfig AppConfig

	appConfigPath := path.Join(appPath, "app.json")
	appConfigBytes, err := os.ReadFile(appConfigPath)
	if err != nil {
		return appConfig, fmt.Errorf("%v\n"+
			"We couldn't find an app.json file on path %q. Maybe try in another using `--path`", err, appPath)
	}
	if err := json.Unmarshal(appConfigBytes, &appConfig); err != nil {
		return appConfig, err
	}

	return appConfig, nil
}
