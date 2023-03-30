## meroxa resources create

Add a resource to your Meroxa resource catalog

### Synopsis

Use the create command to add resources to your Meroxa resource catalog.

```
meroxa resources create [NAME] --type TYPE --url URL [flags]
```

### Examples

```

$ meroxa resource create mybigquery \
    --type bigquery \
    -u "bigquery://$GCP_PROJECT_ID/$GCP_DATASET_NAME" \
    --client-key "$(cat $GCP_SERVICE_ACCOUNT_JSON_FILE)"

$ meroxa resource create sourcedb \
	--type confluentcloud \
	--url kafka+sasl+ssl://$API_KEY:$API_SECRET@<$BOOTSTRAP_SERVER>?sasl_mechanism=plain

$ meroxa resource create meteor \
	--type cosmosdb \
	--url cosmosdb://user:pass@org.documents.azure.com:443/pluto

$ meroxa resource create elasticsearch \
    --type elasticsearch \
    -u "https://$ES_USER:$ES_PASS@$ES_URL:$ES_PORT" \
    --metadata '{"index.prefix": "$ES_INDEX","incrementing.field.name": "$ES_INCREMENTING_FIELD"}'

$ meroxa resource create sourcedb \
	--type kafka \
	--url kafka+sasl+ssl://$KAFKA_USER:$KAFKA_PASS@<$BOOTSTRAP_SERVER>?sasl_mechanism=plain

$ meroxa resource create mongo \
    --type mongodb \
    -u "mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_URL:$MONGO_PORT"

$ meroxa resource create mysqldb \
    --type mysql \
    --url "mysql://$MYSQL_USER:$MYSQL_PASS@$MYSQL_URL:$MYSQL_PORT/$MYSQL_DB"

$ meroxa resource create workspace \
	--type notion \
	--token AbCdEfG123456

$ meroxa resource create workspace \
	--type oracledb \
	--url oracle://user:password@host.com:1521/database

$ meroxa resources create store \
	--type postgres \
	-u "$DATABASE_URL" \
	--metadata '{"logical_replication":"true"}'

$ meroxa resources create warehouse \
	--type redshift \
	-u "$REDSHIFT_URL" \
	--private-key-file ~/.ssh/my-key

$ meroxa resources create datalake \
	--type s3 \
	-u "s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos"

$ meroxa resource create snowflake \
    --type snowflakedb \
    -u "snowflake://$SNOWFLAKE_URL/meroxa_db/stream_data" \
    --username meroxa_user \
    --private-key-file /Users/me/.ssh/snowflake_ed25519

$ meroxa resource create hr \
	--type sqlserver \
	--url "sqlserver://$MSSQL_USER:$MSSQL_PASS@$MSSQL_URL:$MSSQL_PORT/$MSSQL_DB"

$ meroxa resources create slack \
	--type url \
	-u "$WEBHOOK_URL"
```

### Options

```
      --ca-cert string            trusted certificates for verifying resource
      --client-cert string        client certificate for authenticating to the resource
      --client-key string         client private key for authenticating to the resource
      --env string                environment (name or UUID) where resource will be created
  -h, --help                      help for create
  -m, --metadata string           resource metadata
      --password string           password
      --private-key-file string   Path to private key file
      --ssh-private-key string    SSH tunneling private key
      --ssh-url string            SSH tunneling address
      --ssl                       use SSL
      --token string              API Token
      --type string               resource type (required)
  -u, --url string                resource url
      --username string           username
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa resources](meroxa_resources.md)	 - Manage resources on Meroxa

