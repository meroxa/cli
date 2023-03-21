---
createdAt: 
updatedAt: 
title: "meroxa apps describe"
slug: meroxa-apps-describe
url: /cli/cmd/meroxa-apps-describe/
---
## meroxa apps describe

Describe a Turbine Data Application

### Synopsis

This command will fetch details about the Application specified in '--path'
(or current working directory if not specified) on our Meroxa Platform,
or the Application specified by the given name or UUID identifier.

```
meroxa apps describe [NameOrUUID] [--path pwd] [flags]
```

### Examples

```
meroxa apps describe # assumes that the Application is in the current directory
meroxa apps describe --path /my/app
meroxa apps describe NAMEorUUID
```

### Options

```
  -h, --help          help for describe
      --path string   Path to the app directory (default is local directory)
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa apps](/cli/cmd/meroxa-apps/)	 - Manage Turbine Data Applications (Beta)

