---
createdAt: 
updatedAt: 
title: "meroxa secrets update"
slug: meroxa-secrets-update
url: /cli/cmd/meroxa-secrets-update/
---
## meroxa secrets update

Update a Turbine Secret

### Synopsis

This command will update the specified Turbine Secret's data.

```
meroxa secrets update nameOrUUID --data '{"key": "any new data"} [flags]
```

### Examples

```
meroxa secrets update nameOrUUID --data '{"key": "value"}' 
		or 
		meroxa secrets update nameOrUUID 
```

### Options

```
      --data string   Secret's data, passed as a JSON string
  -h, --help          help for update
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

* [meroxa secrets](/cli/cmd/meroxa-secrets/)	 - Manage Turbine Data Applications

