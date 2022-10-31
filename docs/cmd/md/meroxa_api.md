## meroxa api

Invoke Meroxa API

```
meroxa api METHOD PATH [body] [flags]
```

### Examples

```

meroxa api GET /v1/resources
meroxa api POST /v1/resources '{"type":"postgres", "name":"pg", "url":"postgres://.."}'
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

* [meroxa](meroxa.md)	 - The Meroxa CLI

