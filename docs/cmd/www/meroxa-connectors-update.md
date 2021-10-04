---
createdAt: 
updatedAt: 
title: "meroxa connectors update"
slug: meroxa-connectors-update
url: /cli/cmd/meroxa-connectors-update/
---
## meroxa connectors update

Update connector name, configuration or state

```
meroxa connectors update NAME [flags]
```

### Examples

```

meroxa connector update old-name --name new-name' 
meroxa connector update connector-name --state pause' 
meroxa connector update connector-name --config '{"table.name.format":"public.copy"}' 
meroxa connector update connector-name --state restart' 

```

### Options

```
  -c, --config string   new connector configuration
  -h, --help            help for update
      --name string     new connector name
      --state string    new connector state (pause | resume | restart)
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa connectors](/cli/cmd/meroxa-connectors/)	 - Manage connectors on Meroxa

