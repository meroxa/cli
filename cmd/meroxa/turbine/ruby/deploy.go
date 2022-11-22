package turbinerb

import (
	"context"
	"time"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/cmd/meroxa/turbine/ruby/internal"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (t *turbineRbCLI) SetupForDeploy(ctx context.Context) (func(), error) {
	go t.recordServer.Run(ctx)

	cmd := internal.NewTurbineCmd(t.appPath, map[string]string{
		"TURBINE_CORE_SERVER": t.grpcListenAddress,
	})

	if err := utils.RunCMD(ctx, t.logger, cmd); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		t.grpcListenAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	t.recordClient = recordClient{
		ClientConn:           conn,
		TurbineServiceClient: pb.NewTurbineServiceClient(conn),
	}

	return func() {
		t.recordClient.Close()
		t.recordServer.GracefulStop()
	}, nil
}

func (t *turbineRbCLI) NeedsToBuild(ctx context.Context, appName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.recordClient.HasFunctions(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Value, nil
}

func (t *turbineRbCLI) Deploy(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.recordClient.GetSpec(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return string(resp.Spec), nil
}

func (t *turbineRbCLI) GetResources(ctx context.Context, appName string) ([]utils.ApplicationResource, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.recordClient.ListResources(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	var resources []utils.ApplicationResource
	for _, r := range resp.Resources {
		resources = append(resources, utils.ApplicationResource{
			Name: r.Name,
		})
	}
	return resources, nil
}

func (t *turbineRbCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(t.appPath)
}

func (t *turbineRbCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}

func (t *turbineRbCLI) UploadSource(ctx context.Context, appName, appPath, url string) error {
	return utils.UploadSource(ctx, t.logger, utils.Ruby, t.appPath, appName, url)
}
