---
createdAt: 2021-04-15T10:19:35+02:00
updatedAt: 2021-04-15T10:19:35+02:00
title: "meroxa update pipeline"
slug: meroxa-update-pipeline
url: /cli/meroxa-update-pipeline/
---## meroxa update pipeline

Update pipeline state

```
meroxa update pipeline NAME [flags]
```

### Examples

```

meroxa update pipeline old-name --name new-name
meroxa update pipeline pipeline-name --state pause
meroxa update pipeline pipeline-name --metadata '{"key":"value"}'
```

### Options

```
  -h, --help              help for pipeline
  -m, --metadata string   new pipeline metadata
      --name string       new pipeline name
      --state string      new pipeline state
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/meroxa.env)
      --debug           display any debugging information
      --json            output json
```

### SEE ALSO

* [meroxa update](/cli/meroxa-update/)	 - Update a component

