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

meroxa api GET /v1/endpoints
meroxa api POST /v1/endpoints '{"protocol": "HTTP", "stream": "resource-2-499379-public.accounts", "name": "1234"}'
```

### Options

```
  -h, --help   help for api
```

### Options inherited from parent commands

```
      --config string      config file
      --debug              display any debugging information
      --json               output json
      --timeout duration   set the duration of the client timeout in seconds (default 10s) (default 10s)
```

### SEE ALSO

* [meroxa](/cli/cmd/meroxa/)	 - The Meroxa CLI

