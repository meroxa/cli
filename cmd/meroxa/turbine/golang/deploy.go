package turbinego

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/meroxa/turbine-core/pkg/client"

	"github.com/meroxa/cli/cmd/meroxa/turbine"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	turbineGo "github.com/meroxa/turbine-go/deploy"
)

// Deploy runs the binary previously built with the `--deploy` flag which should create all necessary resources.
// TODO: Once all languages are under turbine-core refactor this so it's the same for all languages.
func (t *turbineGoCLI) GetDeploymentSpec(ctx context.Context, imageName, _, _, _, _ string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.bc.GetSpec(ctx, &pb.GetSpecRequest{
		Image: imageName,
	})
	if err != nil {
		return "", err
	}

	return string(resp.Spec), nil
}

// TODO: Once all languages are under turbine-core refactor this so it's the same for all languages.
func (t *turbineGoCLI) GetResources(ctx context.Context) ([]turbine.ApplicationResource, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.bc.ListResources(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	var resources []turbine.ApplicationResource
	for _, r := range resp.Resources {
		resources = append(resources, turbine.ApplicationResource{
			Name:        r.Name,
			Destination: r.Destination,
			Source:      r.Source,
			Collection:  r.Collection,
		})
	}
	return resources, nil
}

// NeedsToBuild reads from the Turbine application to determine whether it needs to be built or not
// this is currently based on the number of functions.
// TODO: Once all languages are under turbine-core refactor this so it's the same for all languages.
func (t *turbineGoCLI) NeedsToBuild(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.bc.HasFunctions(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Value, nil
}

func (t *turbineGoCLI) GetGitSha(ctx context.Context) (string, error) {
	return turbine.GetGitSha(ctx, t.appPath)
}

func (t *turbineGoCLI) GitChecks(ctx context.Context) error {
	return turbine.GitChecks(ctx, t.logger, t.appPath)
}

func (t *turbineGoCLI) CreateDockerfile(_ context.Context, appName, specVersion string) (string, error) {
	return t.appPath, turbineGo.CreateDockerfile(appName, t.appPath, specVersion)
}

func (t *turbineGoCLI) CleanUpBuild(_ context.Context) {
	turbine.CleanupDockerfile(t.logger, t.appPath)
}

func (t *turbineGoCLI) SetupForDeploy(ctx context.Context, gitSha string) (func(), error) {
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	port := rand.Intn(65535-1024) + 1024 //nolint:gomnd
	grpcListenAddress := fmt.Sprintf("localhost:%d", port)

	go t.builder.RunWithAddress(ctx, grpcListenAddress)

	// randomizes the port

	cmd := exec.Command("go", []string{
		"run",
		"./...",
		"build",
		"-gitsha",
		gitSha,
		"-turbine-core-server", grpcListenAddress,
		"-app-path", t.appPath,
	}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = t.appPath
	cmd.Env = os.Environ()

	if err := turbine.RunCMD(ctx, t.logger, cmd); err != nil {
		return nil, err
	}

	c, err := client.DialTimeout(grpcListenAddress, time.Second)
	if err != nil {
		return nil, err
	}
	t.bc = c

	return func() {
		c.Close()
		t.builder.GracefulStop()
	}, nil
}
