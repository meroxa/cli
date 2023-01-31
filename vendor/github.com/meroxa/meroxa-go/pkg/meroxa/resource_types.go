package meroxa

import (
	"context"
	"encoding/json"
	"net/http"
)

type ResourceTypeName string
type ResourceTypeReleaseStage string

const ResourcesTypeBasePath = "/v1/resource-types"
const V2ResourcesTypeBasePath = "/v2/resource-types"

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

	ResourceTypeAirtable                 ResourceTypeName = "airtable"
	ResourceTypeAlgolia                  ResourceTypeName = "algolia"
	ResourceTypeDynamoDB                 ResourceTypeName = "aws_dynamodb"
	ResourceTypeKinesis                  ResourceTypeName = "aws_kinesis"
	ResourceTypeApacheWebLogs            ResourceTypeName = "apache_web_logs"
	ResourceTypeAppDynamics              ResourceTypeName = "app_dynamics"
	ResourceTypeAtlassianConfluence      ResourceTypeName = "atlassian_confluence"
	ResourceTypeAtlassianJira            ResourceTypeName = "atlassian_jira"
	ResourceTypeAzureBlobStorage         ResourceTypeName = "azure_blob_storage"
	ResourceTypeAzureEventHub            ResourceTypeName = "azure_event_hub"
	ResourceTypeBox                      ResourceTypeName = "box"
	ResourceTypeCassandra                ResourceTypeName = "cassandra"
	ResourceTypeClickhouse               ResourceTypeName = "clickhouse"
	ResourceTypeCockroach                ResourceTypeName = "cockroach"
	ResourceTypeDropbox                  ResourceTypeName = "dropbox"
	ResourceTypeFacebookAds              ResourceTypeName = "facebook_ads"
	ResourceTypeFile                     ResourceTypeName = "file"
	ResourceTypeFirebaseFirestore        ResourceTypeName = "firebase_firestore"
	ResourceTypeFirebolt                 ResourceTypeName = "firebolt"
	ResourceTypeFluentbit                ResourceTypeName = "fluentbit"
	ResourceTypeFtpSftp                  ResourceTypeName = "ftp_sftp"
	ResourceTypeGitHub                   ResourceTypeName = "github"
	ResourceTypeGitLab                   ResourceTypeName = "gitlab"
	ResourceTypeGoogleAnalytics          ResourceTypeName = "google_analytics"
	ResourceTypeGoogleCloudStorage       ResourceTypeName = "google_cloud_storage"
	ResourceTypeGoogleDrive              ResourceTypeName = "google_drive"
	ResourceTypeGooglePubSub             ResourceTypeName = "google_pub_sub"
	ResourceTypeGoogleSheets             ResourceTypeName = "google_sheets"
	ResourceTypeHubspot                  ResourceTypeName = "hubspot"
	ResourceTypeIbmDb2                   ResourceTypeName = "ibm_db2"
	ResourceTypeKlayvio                  ResourceTypeName = "klayvio"
	ResourceTypeK8sLogs                  ResourceTypeName = "kubernetes_logs"
	ResourceTypeLogstash                 ResourceTypeName = "logstash"
	ResourceTypeMailchimp                ResourceTypeName = "mailchimp"
	ResourceTypeMarketo                  ResourceTypeName = "marketo"
	ResourceTypeMaterialize              ResourceTypeName = "materialize"
	ResourceTypeMsTeams                  ResourceTypeName = "microsoft_teams"
	ResourceTypeNatsJetstream            ResourceTypeName = "nats_jetstream"
	ResourceTypeNetsuite                 ResourceTypeName = "netsuite"
	ResourceTypeNginx                    ResourceTypeName = "nginx"
	ResourceTypeOpenTelemetry            ResourceTypeName = "open_telemetry"
	ResourceTypeOracle                   ResourceTypeName = "oracle"
	ResourceTypeOsquery                  ResourceTypeName = "osquery"
	ResourceTypePrometheus               ResourceTypeName = "prometheus"
	ResourceTypePulsar                   ResourceTypeName = "pulsar"
	ResourceTypeRedis                    ResourceTypeName = "redis"
	ResourceTypeSalesforceSalesCloud     ResourceTypeName = "salesforce_sales_cloud"
	ResourceTypeSalesforceMarketingCloud ResourceTypeName = "salesforce_marketing_cloud"
	ResourceTypeSalesforcePardot         ResourceTypeName = "salesforce_pardot"
	ResourceTypeSapHana                  ResourceTypeName = "sap_hana"
	ResourceTypeShopify                  ResourceTypeName = "shopify"
	ResourceTypeSlack                    ResourceTypeName = "slack"
	ResourceTypeSocket                   ResourceTypeName = "socket"
	ResourceTypeSplunk                   ResourceTypeName = "splunk"
	ResourceTypeStatsD                   ResourceTypeName = "statsd"
	ResourceTypeStripe                   ResourceTypeName = "stripe"
	ResourceTypeSyslog                   ResourceTypeName = "syslog"
	ResourceTypeTeradata                 ResourceTypeName = "teradata"
	ResourceTypeVitess                   ResourceTypeName = "vitess"
	ResourceTypeWorkday                  ResourceTypeName = "workday"
	ResourceTypeZendeskSupport           ResourceTypeName = "zendesk_support"

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
	CliOnly      bool                     `json:"cli_only"`
}

// ListResourceTypes returns the list of supported resource types.
func (c *client) ListResourceTypes(ctx context.Context) ([]string, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, ResourcesTypeBasePath, nil, nil, nil)
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

// ListResourceTypesV2 returns the list of supported resource types as objects.
func (c *client) ListResourceTypesV2(ctx context.Context) ([]ResourceType, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, V2ResourcesTypeBasePath, nil, nil, nil)
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
