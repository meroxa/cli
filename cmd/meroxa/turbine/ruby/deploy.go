package turbinerb

import (
	"context"

	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
)

func (t *turbineRbCLI) NeedsToBuild(ctx context.Context, appName string) (bool, error) {

	return false, nil
}

func (t *turbineRbCLI) Deploy(ctx context.Context, imageName, appName, gitSha, specVersion, accountUUID string) (string, error) {
	var (
		deploymentSpec string
		err            error
	)
	return deploymentSpec, err
}

func (t *turbineRbCLI) GetResources(ctx context.Context, appName string) ([]utils.ApplicationResource, error) {
	var resources []utils.ApplicationResource
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
