---
createdAt: 
updatedAt: 
title: "meroxa pipelines update"
slug: meroxa-pipelines-update
url: /cli/cmd/meroxa-pipelines-update/
---
## meroxa pipelines update

Update pipeline name, state or metadata

```
meroxa pipelines update NAME [flags]
```

### Examples

```

meroxa pipeline update old-name --name new-name
meroxa pipeline update pipeline-name --state pause
meroxa pipeline update pipeline-name --metadata '{"key":"value"}'
```

### Options

```
  -h, --help              help for update
  -m, --metadata string   new pipeline metadata
      --name string       new pipeline name
      --state string      new pipeline state
```

### Options inherited from parent commands

```
      --config string      config file
      --debug              display any debugging information
      --json               output json
      --timeout duration   set the duration of the client timeout in seconds (default 10s) (default 10s)
```

### SEE ALSO

* [meroxa pipelines](/cli/cmd/meroxa-pipelines/)	 - Manage pipelines on Meroxa

