package frameworks

import (
	"context"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
)

type CLI interface {
	Build(ctx context.Context, appName string, platform bool) error
	CleanUpBinaries(ctx context.Context, appName string)
	Deploy(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error)
	GetResources(ctx context.Context, appName string) ([]turbine.ApplicationResource, error)
	GetVersion(ctx context.Context) (string, error)
	GitInit(ctx context.Context, name string) error
	GitChecks(ctx context.Context) error
	GetGitSha(ctx context.Context) (string, error)
	Init(ctx context.Context, name string) error
	NeedsToBuild(ctx context.Context, appName string) (bool, error)
	Run(ctx context.Context) error
	Upgrade(vendor bool) error
	CreateDockerfile(ctx context.Context, appName string) (string, error)
	CleanUpBuild(ctx context.Context)
	SetupForDeploy(ctx context.Context, gitSha string) (func(), error)
}
