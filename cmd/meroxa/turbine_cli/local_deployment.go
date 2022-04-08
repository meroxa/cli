package turbinecli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/pkg/archive"
	turbine "github.com/meroxa/turbine-go/deploy"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/meroxa/cli/log"
)

type LocalDeploy struct {
	DockerHubUserNameEnv    string
	DockerHubAccessTokenEnv string
	Enabled                 bool
	TempPath                string
	Lang                    string
}

func (ld *LocalDeploy) getAuthConfig() string {
	dhUsername := ld.DockerHubUserNameEnv
	dhAccessToken := ld.DockerHubAccessTokenEnv
	authConfig := types.AuthConfig{
		Username:      dhUsername,
		Password:      dhAccessToken,
		ServerAddress: "https://index.docker.io/v1/",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	return base64.URLEncoding.EncodeToString(authConfigBytes)
}

// GetDockerImageName Will create the image via DockerHub.
func (ld *LocalDeploy) GetDockerImageName(ctx context.Context, l log.Logger, appPath, appName, lang string) (string, error) {
	l.Info(ctx, "\t  Using DockerHub...")
	// fqImageName will be eventually taken from the build endpoint.
	fqImageName := ld.DockerHubUserNameEnv + "/" + appName

	err := ld.buildImage(ctx, l, appPath, fqImageName)
	if err != nil {
		return "", err
	}

	// this will go away when using the new service.
	err = ld.pushImage(l, fqImageName)
	if err != nil {
		l.Errorf(ctx, "\t 𐄂 Unable to push image %q", fqImageName)
		return "", err
	}

	l.Infof(ctx, "\t✔ Image built!")
	return fqImageName, nil
}

func (ld *LocalDeploy) buildImage(ctx context.Context, l log.Logger, pwd, imageName string) error {
	l.Infof(ctx, "\t  Building image %q located at %q", imageName, pwd)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		l.Errorf(ctx, "\t 𐄂 Unable to init docker client")
		return err
	}

	// Generate dockerfile
	if ld.Lang == "golang" {
		err = turbine.CreateDockerfile(pwd)
		if err != nil {
			return err
		}
	}

	if ld.Lang == "javascript" || ld.Lang == "python" {
		pwd = ld.TempPath
	}
	// Read local Dockerfile
	tar, err := archive.TarWithOptions(pwd, &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{".git", "fixtures"},
	})
	if err != nil {
		l.Errorf(ctx, "\t 𐄂 Unable to create tar")
		return err
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
		l.Errorf(ctx, "\t 𐄂 Unable to build docker image")
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			l.Errorf(ctx, "\t 𐄂 Unable to close docker build response body; %s", err)
		}
	}(resp.Body)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		l.Errorf(ctx, "\t 𐄂 Unable to read image build response")
		return err
	}
	dockerBuildOutput := buf.String()

	re := regexp.MustCompile(`{"errorDetail":{"message":"([^"]+)"}`)
	matches := re.FindAllStringSubmatch(dockerBuildOutput, -1)
	if len(matches) != 0 {
		const subMatchArraySize = 2
		errMsg := ""
		for _, str := range matches {
			if len(str) < subMatchArraySize {
				continue
			}
			errMsg += "\n" + str[1]
		}
		err = fmt.Errorf("%s", errMsg)
		l.Errorf(ctx, "\t 𐄂 Unable to build docker image")
		return err
	}
	l.Info(ctx, dockerBuildOutput)

	if ld.Lang == "golang" {
		// Cleanup
		return os.Remove(filepath.Join(pwd, "Dockerfile"))
	}
	return nil
}

func (ld *LocalDeploy) pushImage(l log.Logger, imageName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	authConfig := ld.getAuthConfig()
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
