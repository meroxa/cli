//go:generate mockgen -source=interface.go -package=mock -destination=mock/cli_mock.go CLI

package turbine

import (
	"context"

	"github.com/meroxa/cli/log"
)

type CLI interface {
	GitInit(context.Context, string) error
	CheckUncommittedChanges(context.Context, log.Logger, string) error
	CleanupDockerfile(log.Logger, string)
	GetGitSha(context.Context, string) (string, error)
	GetDeploymentSpec(context.Context, string) (string, error)
	GetResources(context.Context) ([]ApplicationResource, error)
	NeedsToBuild(context.Context) (bool, error)
	GetVersion(context.Context) (string, error)
	Init(context.Context, string) error
	Run(context.Context) error
	CreateDockerfile(context.Context, string) (string, error)
	StartGrpcServer(context.Context, string) (func(), error)
}
