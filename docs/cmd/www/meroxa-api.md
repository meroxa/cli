---
createdAt: 
updatedAt: 
title: "meroxa api"
slug: meroxa-api
url: /cli/cmd/meroxa-api/
---
## meroxa api

Invoke Meroxa API

```
meroxa api METHOD PATH [body] [flags]
```

### Examples

```

meroxa api GET /api/collectoions/{collection}/records
meroxa api POST /api/collectoions/{collection}/records '{"type":"postgres", "name":"pg", "url":"postgres://.."}'
```

### Options

```
  -h, --help   help for api
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa](/docs/cmd/www/meroxa.md)	 - The Meroxa CLI

