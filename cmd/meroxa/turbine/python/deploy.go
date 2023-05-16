package turbinepy

import (
	"context"
	"time"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/python/internal"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

// NeedsToBuild determines if the app has functions or not.

func (t *turbinePyCLI) NeedsToBuild(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.bc.HasFunctions(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Value, nil
}

func (t *turbinePyCLI) CreateDockerfile(ctx context.Context, _ string) (string, error) {
	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandBuild,
		map[string]string{},
		t.appPath,
	)
	return t.appPath, utils.RunCMD(ctx, t.logger, cmd)
}

func (t *turbinePyCLI) CleanUpBuild(_ context.Context) {
	utils.CleanupDockerfile(t.logger, t.appPath)
}

func (t *turbinePyCLI) SetupForDeploy(ctx context.Context, gitSha string) (func(), error) {
	go t.builder.RunAddr(ctx, t.grpcListenAddress)

	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandRecord,
		map[string]string{
			"TURBINE_CORE_SERVER": t.grpcListenAddress,
		},
		t.appPath,
		gitSha,
	)

	if err := utils.RunCMD(ctx, t.logger, cmd); err != nil {
		return nil, err
	}

	c, err := client.DialTimeout(t.grpcListenAddress, time.Second)
	if err != nil {
		return nil, err
	}
	t.bc = c

	return func() {
		c.Close()
		t.builder.GracefulStop()
	}, nil
}

// Deploy creates Application entities.
func (t *turbinePyCLI) Deploy(ctx context.Context, imageName, _, _, _, _ string) (string, error) {
	resp, err := t.bc.GetSpec(ctx, &pb.GetSpecRequest{
		Image: imageName,
	})
	if err != nil {
		return "", err
	}
	return string(resp.Spec), nil
}

func (t *turbinePyCLI) GetDeploymentSpec(ctx context.Context, imageName, _, _, _, _ string) (string, error) {
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

// GetResources asks turbine for a list of resources used by the given app.
func (t *turbinePyCLI) GetResources(ctx context.Context) ([]utils.ApplicationResource, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.bc.ListResources(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	var resources []utils.ApplicationResource
	for _, r := range resp.Resources {
		resources = append(resources, utils.ApplicationResource{
			Name:        r.Name,
			Destination: r.Destination,
			Source:      r.Source,
			Collection:  r.Collection,
		})
	}
	return resources, nil
}
