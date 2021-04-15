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
      --config string   config file (default is $HOME/meroxa.env)
      --debug           display any debugging information
      --json            output json
```

### SEE ALSO

* [meroxa](meroxa.md)	 - The Meroxa CLI

