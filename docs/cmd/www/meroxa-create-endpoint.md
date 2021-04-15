---
createdAt: 2021-04-15T17:48:51+02:00
updatedAt: 2021-04-15T17:48:51+02:00
title: "meroxa create endpoint"
slug: meroxa-create-endpoint
url: /cli/meroxa-create-endpoint/
---
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
      --config string      config file (default is $HOME/meroxa.env)
      --debug              display any debugging information
      --json               output json
      --timeout duration   set the client timeout (default 10s)
```

### SEE ALSO

* [meroxa create](/cli/meroxa-create/)	 - Create Meroxa pipeline components

