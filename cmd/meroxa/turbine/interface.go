//go:generate mockgen -source=interface.go -package=mock -destination=mock/cli_mock.go CLI

package turbine

import (
	"context"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/grpc"
)

type CLI interface {
	GetDeploymentSpec(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error)
	GetResources(ctx context.Context) ([]ApplicationResource, error)
	GetVersion(ctx context.Context) (string, error)
	GitInit(ctx context.Context, name string) error
	GitChecks(ctx context.Context) error
	GetGitSha(ctx context.Context) (string, error)
	Init(ctx context.Context, name string) error
	NeedsToBuild(ctx context.Context) (bool, error)
	Run(ctx context.Context) error
	Upgrade(vendor bool) error
	CreateDockerfile(ctx context.Context, appName, specVersion string) (string, error)
	CleanUpBuild(ctx context.Context)
	SetupForDeploy(ctx context.Context, gitSha string) (func(), error)
}

type Server interface {
	Run(context.Context)
	GracefulStop()
}

type Core struct {
	Builder           Server
	SpecBuilderClient SpecBuilderClient
	Runner            Server
	GrpcListenAddress string
}

type SpecBuilderClient struct {
	*grpc.ClientConn
	pb.TurbineServiceClient
}
