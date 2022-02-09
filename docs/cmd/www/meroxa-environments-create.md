---
createdAt: 
updatedAt: 
title: "meroxa environments create"
slug: meroxa-environments-create
url: /cli/cmd/meroxa-environments-create/
---
## meroxa environments create

Create an environment

```
meroxa environments create NAME [flags]
```

### Examples

```

meroxa env create my-env --type self_hosted --provider aws --region us-east-1 --config aws_access_key_id=my_access_key --config aws_secret_access_key=my_access_secret

```

### Options

```
  -c, --config strings    environment configuration based on type and provider (e.g.: --config aws_access_key_id=my_access_key --config aws_secret_access_key=my_access_secret)
  -h, --help              help for create
      --provider string   environment cloud provider to use
      --region string     environment region
      --type string       environment type, when not specified
  -y, --yes               skip confirmation prompt
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

