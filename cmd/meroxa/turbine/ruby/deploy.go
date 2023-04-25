package turbinerb

import (
	"context"
	"time"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (t *turbineRbCLI) CleanUpBuild(_ context.Context) {
	utils.CleanupDockerfile(t.logger, t.appPath)
}

func (t *turbineRbCLI) CreateDockerfile(ctx context.Context, _ string) (string, error) {
	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandBuild,
		map[string]string{},
	)
	return t.appPath, utils.RunCMD(ctx, t.logger, cmd)
}

func (t *turbineRbCLI) SetupForDeploy(ctx context.Context, gitSha string) (func(), error) {
	go t.builder.Run(ctx)

	cmd := internal.NewTurbineCmd(t.appPath,
		internal.TurbineCommandRecord,
		map[string]string{
			"TURBINE_CORE_SERVER": t.grpcListenAddress,
		},
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

func (t *turbineRbCLI) NeedsToBuild(ctx context.Context, _ string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.bc.HasFunctions(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Value, nil
}

func (t *turbineRbCLI) Deploy(ctx context.Context, imageName, _, _, _, _ string) (string, error) {
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

func (t *turbineRbCLI) GetResources(ctx context.Context, _ string) ([]utils.ApplicationResource, error) {
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

func (t *turbineRbCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(ctx, t.appPath)
}

func (t *turbineRbCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}
