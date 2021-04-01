## meroxa create resource

Add a resource to your Meroxa resource catalog

### Synopsis

Use the add command to add resources to your Meroxa resource catalog.

```
meroxa create resource <resource-name> --type <resource-type> [flags]
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
      --from string       resource name to use as source
  -h, --help              help for resource
      --input string      command delimited list of input streams
      --pipeline string   ID of pipeline to attach connector to
      --to string         resource name to use as destination
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/meroxa.env)
      --json            output json
```

### SEE ALSO

* [meroxa create](meroxa_create.md)	 - Create Meroxa pipeline components

