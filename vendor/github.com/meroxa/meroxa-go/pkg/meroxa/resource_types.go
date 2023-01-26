package meroxa

import (
	"context"
	"encoding/json"
	"net/http"
)

type ResourceTypeName string
type ResourceTypeReleaseStage string

const ResourceTypeFormConfigHumanReadableKey = "label"

const (
	ResourceTypePostgres       ResourceTypeName = "postgres"
	ResourceTypeMysql          ResourceTypeName = "mysql"
	ResourceTypeRedshift       ResourceTypeName = "redshift"
	ResourceTypeUrl            ResourceTypeName = "url"
	ResourceTypeS3             ResourceTypeName = "s3"
	ResourceTypeMongodb        ResourceTypeName = "mongodb"
	ResourceTypeElasticsearch  ResourceTypeName = "elasticsearch"
	ResourceTypeSnowflake      ResourceTypeName = "snowflakedb"
	ResourceTypeBigquery       ResourceTypeName = "bigquery"
	ResourceTypeSqlserver      ResourceTypeName = "sqlserver"
	ResourceTypeCosmosdb       ResourceTypeName = "cosmosdb"
	ResourceTypeKafka          ResourceTypeName = "kafka"
	ResourceTypeConfluentCloud ResourceTypeName = "confluentcloud"
	ResourceTypeNotion         ResourceTypeName = "notion"

	ResourceTypeReleaseStageGA         ResourceTypeReleaseStage = "ga"
	ResourceTypeReleaseStageBeta       ResourceTypeReleaseStage = "beta"
	ResourceTypeReleaseStageDevPreview ResourceTypeReleaseStage = "developer_preview"
)

type ResourceType struct {
	UUID         string                   `json:"uuid"`
	Name         string                   `json:"name"`
	ReleaseStage ResourceTypeReleaseStage `json:"release_stage"`
	Categories   []string                 `json:"categories"`
	FormConfig   map[string]interface{}   `json:"form_config"`
	OptedIn      bool                     `json:"opted_in"`
	HasAccess    bool                     `json:"has_access"`
	CLIOnly      bool                     `json:"cli_only"`
}

// ListResourceTypes returns the list of supported resources
func (c *client) ListResourceTypes(ctx context.Context) ([]ResourceType, error) {
	path := "/v1/resource-types"

	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var supportedTypes []ResourceType
	err = json.NewDecoder(resp.Body).Decode(&supportedTypes)
	if err != nil {
		return nil, err
	}

	return supportedTypes, nil
}
