---
createdAt: 
updatedAt: 
title: "meroxa resources create"
slug: meroxa-resources-create
url: /cli/cmd/meroxa-resources-create/
---
## meroxa resources create

Add a resource to your Meroxa resource catalog

### Synopsis

Use the create command to add resources to your Meroxa resource catalog.

```
meroxa resources create [NAME] --type TYPE --url URL [flags]
```

### Examples

```

$ meroxa resources create store --type postgres -u "$DATABASE_URL" --metadata '{"logical_replication":"true"}'
$ meroxa resources create datalake --type s3 -u "s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos"
$ meroxa resources create warehouse --type redshift -u "$REDSHIFT_URL"
$ meroxa resources create slack --type url -u "$WEBHOOK_URL"
$ meroxa resource create mysqldb \
    --type mysql \
    --url "mysql://$MYSQL_USER:$MYSQL_PASS@$MYSQL_URL:$MYSQL_PORT/$MYSQL_DB"

$ meroxa resource create mongo \
    --type mongodb \
    -u "mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_URL:$MONGO_PORT"

$ meroxa resource create elasticsearch \
    --type elasticsearch \
    -u "https://$ES_USER:$ES_PASS@$ES_URL:$ES_PORT" \
    --metadata '{"index.prefix": "$ES_INDEX","incrementing.field.name": "$ES_INCREMENTING_FIELD"}'

$ meroxa resource create mybigquery \
    --type bigquery \
    -u "bigquery://$GCP_PROJECT_ID/$GCP_DATASET_NAME" \
    --client-key "$(cat $GCP_SERVICE_ACCOUNT_JSON_FILE)"

$ meroxa resource create snowflake \
    --type snowflakedb \
    -u "snowflake://$SNOWFLAKE_URL/meroxa_db/stream_data" \
    --username meroxa_user \
    --password $SNOWFLAKE_PRIVATE_KEY
```

### Options

```
      --ca-cert string           trusted certificates for verifying resource
      --client-cert string       client certificate for authenticating to the resource
      --client-key string        client private key for authenticating to the resource
      --env string               environment (name or UUID) where resource will be created
  -h, --help                     help for create
  -m, --metadata string          resource metadata
      --password string          password
      --ssh-private-key string   SSH tunneling private key
      --ssh-url string           SSH tunneling address
      --ssl                      use SSL
      --type string              resource type (required)
  -u, --url string               resource url (required)
      --username string          username
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa resources](/cli/cmd/meroxa-resources/)	 - Manage resources on Meroxa

