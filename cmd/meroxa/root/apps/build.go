package apps

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/meroxa/cli/log"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildGoApp(ctx context.Context, l log.Logger, path string, platform bool) error {
	var cmd *exec.Cmd
	if platform {
		cmd = exec.Command("go", "build", "--tags", "platform", "-o", "./"+path, "./"+path+"/...")
	} else {
		cmd = exec.Command("go", "build", path+"/...")
	}

	l.Info(ctx, "building app...\n")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		l.Error(ctx, string(stdout))
		return err
	}
	// TODO: parse output for build errors
	l.Info(ctx, "build complete!")
	return nil
}

func deployApp(ctx context.Context, l log.Logger, projPath, projName string, imageName string) error {
	l.Info(ctx, "deploying app...\n")
	cmd := exec.Command(projPath+"/"+projName, "--deploy", "--imagename", imageName)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	l.Info(ctx, string(stdout))
	l.Info(ctx, "deploy complete!")
	return nil
}

func buildImage(ctx context.Context, l log.Logger, path string, name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		l.Errorf(ctx, "unable to init docker client; %s", err)
	}
	// Read local Dockerfile
	tar, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{"simple", ".git", "fixtures"},
	})
	if err != nil {
		l.Errorf(ctx, "unable to create tar; %s", err)
	}

	buildOptions := types.ImageBuildOptions{
		Context:    tar,
		Dockerfile: "Dockerfile",
		Remove:     true,
		Tags:       []string{name}}

	resp, err := cli.ImageBuild(
		ctx,
		tar,
		buildOptions,
	)
	if err != nil {
		l.Errorf(ctx, "unable to build docker image; %s", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		l.Errorf(ctx, "unable to read image build response; %s", err)
	}
	return nil
}

func pushImage(ctx context.Context, l log.Logger, imageName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	authConfig := getAuthConfig()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	l.Infof(ctx, "pushing image %s to container registry", imageName)
	opts := types.ImagePushOptions{RegistryAuth: authConfig}
	rd, err := cli.ImagePush(ctx, imageName, opts)
	if err != nil {
		return err
	}
	defer rd.Close()

	_, err = io.Copy(os.Stdout, rd)
	if err != nil {
		return err
	}
	l.Info(ctx, "image successfully pushed to container registry!")

	return nil
}

func getAuthConfig() string {
	dhUsername := os.Getenv("DOCKER_HUB_USERNAME")
	dhPassword := os.Getenv("DOCKER_HUB_PASSWORD")
	authConfig := types.AuthConfig{
		Username:      dhUsername,
		Password:      dhPassword,
		ServerAddress: "https://index.docker.io/v1/",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	return base64.URLEncoding.EncodeToString(authConfigBytes)
}

func formatName(accountName, imageName string) string {
	return strings.Join([]string{accountName, imageName}, "/")
}

func prependAccount(imageName string) string {
	account := os.Getenv("DOCKER_HUB_USERNAME")
	return strings.Join([]string{account, imageName}, "/")
}
