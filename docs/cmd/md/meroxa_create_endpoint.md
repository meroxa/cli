## meroxa create endpoint

Create an endpoint

### Synopsis

Use create endpoint to expose an endpoint to a connector stream

```
meroxa create endpoint [NAME] [flags]
```

### Examples

```

meroxa create endpoint my-endpoint --protocol http --stream my-stream
```

### Options

```
  -h, --help              help for endpoint
  -p, --protocol string   protocol, value can be http or grpc (required)
  -s, --stream string     stream name (required)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/meroxa.env)
      --debug           display any debugging information
      --json            output json
```

### SEE ALSO

* [meroxa create](meroxa_create.md)	 - Create Meroxa pipeline components

