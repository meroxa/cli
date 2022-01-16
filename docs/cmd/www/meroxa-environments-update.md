---
createdAt: 
updatedAt: 
title: "meroxa environments update"
slug: meroxa-environments-update
url: /cli/cmd/meroxa-environments-update/
---
## meroxa environments update

Update an environment

```
meroxa environments update NAMEorUUID [flags]
```

### Examples

```

meroxa env update my-env --name new-name --config aws_access_key_id=my_access_key --config aws_access_secret=my_access_secret"

```

### Options

```
  -c, --config strings   updated environment configuration based on type and provider (e.g.: --config aws_access_key_id=my_access_key --config aws_access_secret=my_access_secret)
  -h, --help             help for update
      --name string      updated environment name, when specified
  -y, --yes              skip confirmation prompt
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa environments](/cli/cmd/meroxa-environments/)	 - Manage environments on Meroxa

