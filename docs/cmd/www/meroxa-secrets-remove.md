---
createdAt: 
updatedAt: 
title: "meroxa secrets remove"
slug: meroxa-secrets-remove
url: /cli/cmd/meroxa-secrets-remove/
---
## meroxa secrets remove

Remove a Conduit Secret

### Synopsis

This command will remove the secret specified either by name or id

```
meroxa secrets remove [--path pwd] [flags]
```

### Examples

```
meroxa apps remove nameOrUUID
```

### Options

```
  -f, --force   skip confirmation
  -h, --help    help for remove
```

### Options inherited from parent commands

```
      --api-url string           API url
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa secrets](/cli/cmd/meroxa-secrets/)	 - Manage Conduit Data Applications

