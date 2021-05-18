## meroxa resources create

Add a resource to your Meroxa resource catalog

### Synopsis

Use the create command to add resources to your Meroxa resource catalog.

```
meroxa resources create [NAME] --type TYPE --url URL [flags]
```

### Examples

```

meroxa resources create store --type postgres -u $DATABASE_URL --metadata '{"logical_replication":true}'
meroxa resources create datalake --type s3 -u "s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos"
meroxa resources create warehouse --type redshift -u $REDSHIFT_URL
meroxa resources create slack --type url -u $WEBHOOK_URL

```

### Options

```
      --ca-cert string       trusted certificates for verifying resource
      --client-cert string   client certificate for authenticating to the resource
      --client-key string    client private key for authenticating to the resource
  -h, --help                 help for create
  -m, --metadata string      resource metadata
      --password string      password
      --ssl                  use SSL
      --type string          resource type
  -u, --url string           resource url
      --username string      username
```

### Options inherited from parent commands

```
      --config string      config file
      --debug              display any debugging information
      --json               output json
      --timeout duration   set the client timeout (default 10s)
```

### SEE ALSO

* [meroxa resources](meroxa_resources.md)	 - Manage resources on Meroxa

