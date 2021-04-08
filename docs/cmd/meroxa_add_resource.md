## meroxa add resource

Add a resource to your Meroxa resource catalog

### Synopsis

Use the add command to add resources to your Meroxa resource catalog.

```
meroxa add resource [NAME] --type TYPE [flags]
```

### Examples

```

meroxa add resource store --type postgres -u $DATABASE_URL
meroxa add resource datalake --type s3 -u "s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos"
meroxa add resource warehouse --type redshift -u $REDSHIFT_URL
meroxa add resource slack --type url -u $WEBHOOK_URL

```

### Options

```
      --credentials string   resource credentials
  -h, --help                 help for resource
  -m, --metadata string      resource metadata
      --type string          resource type
  -u, --url string           resource url
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/meroxa.env)
      --debug           display any debugging information
      --json            output json
```

### SEE ALSO

* [meroxa add](meroxa_add.md)	 - Add a resource to your Meroxa resource catalog

