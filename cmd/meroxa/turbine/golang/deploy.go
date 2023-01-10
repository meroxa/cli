package turbinego

import (
	"context"
	"time"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	turbineGo "github.com/meroxa/turbine-go/deploy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Deploy runs the binary previously built with the `--deploy` flag which should create all necessary resources.
func (t *turbineGoCLI) Deploy(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.recordClient.GetSpec(ctx, &pb.GetSpecRequest{
		Image: imageName,
	})
	if err != nil {
		return "", err
	}

	return string(resp.Spec), nil
}

func (t *turbineGoCLI) GetResources(ctx context.Context, appName string) ([]utils.ApplicationResource, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.recordClient.ListResources(ctx, &emptypb.Empty{})
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

// NeedsToBuild reads from the Turbine application to determine whether it needs to be built or not
// this is currently based on the number of functions.
func (t *turbineGoCLI) NeedsToBuild(ctx context.Context, appName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := t.recordClient.HasFunctions(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Value, nil
}

func (t *turbineGoCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(t.appPath)
}

func (t *turbineGoCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}

func (t *turbineGoCLI) CreateDockerfile(ctx context.Context, appName string) (string, error) {
	return t.appPath, turbineGo.CreateDockerfile(appName, t.appPath)
}

func (t *turbineGoCLI) CleanUpBuild(ctx context.Context) {
	utils.CleanupDockerfile(t.logger, t.appPath)
}

func (t *turbineGoCLI) SetupForDeploy(ctx context.Context, gitSha string) (func(), error) {
	go t.recordServer.Run(ctx)

	cmd := NewTurbineGoCmd(t.appPath, t.appName, TurbineCommandRun, map[string]string{
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

	t.recordClient = utils.RecordClient{
		ClientConn:           conn,
		TurbineServiceClient: pb.NewTurbineServiceClient(conn),
	}

	return func() {
		t.recordClient.Close()
		t.recordServer.GracefulStop()
	}, nil

}
