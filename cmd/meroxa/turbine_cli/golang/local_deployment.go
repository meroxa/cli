// Package turbineGo TODO: Reorganize this in a different pkg
package turbinego

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/docker/docker/pkg/archive"
	turbine "github.com/meroxa/turbine/deploy"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/meroxa/cli/log"
)

func (gd *Deploy) getAuthConfig() string {
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

// getDockerImageName Will create the image via DockerHub.
func (gd *Deploy) getDockerImageName(ctx context.Context, l log.Logger, appPath, appName string) (string, error) {
	// fqImageName will be eventually taken from the build endpoint.
	fqImageName := gd.DockerHubUserNameEnv + "/" + appName

	err := gd.buildImage(ctx, l, appPath, fqImageName)
	if err != nil {
		l.Errorf(ctx, "unable to build image; %q\n%s", fqImageName, err)
		return "", err
	}

	// this will go away when using the new service.
	err = gd.pushImage(l, fqImageName)
	if err != nil {
		l.Errorf(ctx, "unable to push image; %q\n%s", fqImageName, err)
		return "", err
	}

	return fqImageName, nil
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

func (gd *Deploy) pushImage(l log.Logger, imageName string) error {
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
