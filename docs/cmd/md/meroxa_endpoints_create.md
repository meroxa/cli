## meroxa endpoints create

Create an endpoint

### Synopsis

Use create endpoint to expose an endpoint to a connector stream

```
meroxa endpoints create [NAME] [flags]
```

### Examples

```
meroxa endpoints create my-endpoint --protocol http --stream my-stream
```

### Options

```
  -h, --help              help for create
  -p, --protocol string   protocol, value can be http or grpc
  -s, --stream string     stream name
```

### Options inherited from parent commands

```
      --config string      config file
      --debug              display any debugging information
      --json               output json
      --timeout duration   set the client timeout (default 10s)
```

### SEE ALSO

* [meroxa endpoints](meroxa_endpoints.md)	 - Manage endpoints on Meroxa

