---
createdAt: 2021-04-16T19:26:31+02:00
updatedAt: 2021-04-16T19:26:31+02:00
title: "meroxa api"
slug: meroxa-api
url: /cli/meroxa-api/
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
      --config string      config file (default is $HOME/meroxa.env)
      --debug              display any debugging information
      --json               output json
      --timeout duration   set the client timeout (default 10s)
```

### SEE ALSO

* [meroxa](/cli/meroxa/)	 - The Meroxa CLI

