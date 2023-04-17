package flink

import (
	"context"
	"encoding/json"
	"fmt"
	utils "github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func (t *flinkCLI) Deploy(ctx context.Context, jarURL, appName, gitSha, specVersion, accountUUID string) (string, error) {
	var spec ir.DeploymentSpec

	fc, err := utils.ReadFlinkConfigFile(t.appPath)
	if err != nil {
		return "", err
	}

	ver, _ := t.GetVersion(ctx)
	spec.Definition.GitSha = gitSha
	spec.Definition.Metadata = ir.MetadataSpec{
		Turbine: ir.TurbineSpec{
			Language: utils.Java,
			Version:  ver,
		},
		SpecVersion: specVersion}
	spec.Connectors = []ir.ConnectorSpec{
		{
			Type:       ir.ConnectorSource,
			Resource:   fc.SourceConfig.Name,
			Collection: fc.SourceConfig.Collection,
			Config:     fc.SourceConfig.ConnectorConfig,
		},
		{},
	}
	spec.FlinkJobs = []ir.FlinkJobSpec{
		{
			Name:   appName,
			JarURL: jarURL,
		},
	}

	bytes, err := json.Marshal(spec)
	return string(bytes), err
}

func (t *flinkCLI) GetResources(_ context.Context, _ string) ([]utils.ApplicationResource, error) {
	var resources []utils.ApplicationResource

	fc, err := utils.ReadFlinkConfigFile(t.appPath)
	if err != nil {
		return resources, err
	}
	if fc.SourceConfig == nil || fc.DestinationConfig == nil {
		return resources, fmt.Errorf("config.json not fully populated correctly")
	}

	resources = append(
		resources,
		utils.ApplicationResource{
			Name:       fc.SourceConfig.Name,
			Source:     true,
			Collection: fc.SourceConfig.Collection,
		})
	resources = append(
		resources,
		utils.ApplicationResource{
			Name:        fc.DestinationConfig.Name,
			Destination: true,
			Collection:  fc.DestinationConfig.Collection,
		})

	return resources, nil
}

// NeedsToBuild is not relevant to flink jobs.
func (t *flinkCLI) NeedsToBuild(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func (t *flinkCLI) GetGitSha(ctx context.Context) (string, error) {
	return utils.GetGitSha(ctx, t.appPath)
}

func (t *flinkCLI) GitChecks(ctx context.Context) error {
	return utils.GitChecks(ctx, t.logger, t.appPath)
}

func (t *flinkCLI) CreateDockerfile(_ context.Context, _ string) (string, error) {
	return t.appPath, nil
}

func (t *flinkCLI) CleanUpBuild(_ context.Context) {
}

func (t *flinkCLI) SetupForDeploy(_ context.Context, _ string) (func(), error) {
	return func() {}, nil
}
