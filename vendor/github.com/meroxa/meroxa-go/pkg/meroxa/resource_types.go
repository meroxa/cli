package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResourceType string

const (
	ResourceTypePostgres      ResourceType = "postgres"
	ResourceTypeMysql         ResourceType = "mysql"
	ResourceTypeRedshift      ResourceType = "redshift"
	ResourceTypeUrl           ResourceType = "url"
	ResourceTypeS3            ResourceType = "s3"
	ResourceTypeMongodb       ResourceType = "mongodb"
	ResourceTypeElasticsearch ResourceType = "elasticsearch"
	ResourceTypeSnowflake     ResourceType = "snowflakedb"
	ResourceTypeBigquery      ResourceType = "bigquery"
	ResourceTypeSqlserver     ResourceType = "sqlserver"
	ResourceTypeCosmosdb      ResourceType = "cosmosdb"
)

// ListResourceTypes returns the list of supported resources
func (c *client) ListResourceTypes(ctx context.Context) ([]string, error) {
	path := fmt.Sprintf("/v1/resource-types")

	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var supportedTypes []string
	err = json.NewDecoder(resp.Body).Decode(&supportedTypes)
	if err != nil {
		return nil, err
	}

	return supportedTypes, nil
}
