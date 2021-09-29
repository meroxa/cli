---
createdAt: 
updatedAt: 
title: "meroxa endpoints create"
slug: meroxa-endpoints-create
url: /cli/cmd/meroxa-endpoints-create/
---
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
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s) (default 10s)
```

### SEE ALSO

* [meroxa endpoints](/cli/cmd/meroxa-endpoints/)	 - Manage endpoints on Meroxa

