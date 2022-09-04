package turbine

import (
	"context"
)

type CLI interface {
	Build(ctx context.Context, appName string, platform bool) (string, error)
	CleanUpBinaries(ctx context.Context, appName string)
	Deploy(ctx context.Context, imageName, appName, gitSha, specVersion string) error
	GetResources(ctx context.Context, appName string) ([]ApplicationResource, error)
	GitInit(ctx context.Context, name string) error
	GitChecks(ctx context.Context) error
	GetGitSha(ctx context.Context) (string, error)
	Init(ctx context.Context, name string) error
	NeedsToBuild(ctx context.Context, appName string) (bool, error)
	Run(ctx context.Context) error
	Upgrade(vendor bool) error
	UploadSource(ctx context.Context, appName, tempPath, url string) error
}
