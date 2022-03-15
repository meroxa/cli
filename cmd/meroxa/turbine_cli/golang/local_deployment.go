// Package turbineGo TODO: Reorganize this in a different pkg
package turbineGo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"time"

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
