---
createdAt: 
updatedAt: 
title: "meroxa apps deploy"
slug: meroxa-apps-deploy
url: /cli/cmd/meroxa-apps-deploy/
---
## meroxa apps deploy

Deploy a Turbine Data Application (Beta)

### Synopsis

This command will deploy the application specified in '--path'
(or current working directory if not specified) to our Meroxa Platform.
If deployment was successful, you should expect an application you'll be able to fully manage


```
meroxa apps deploy [flags]
```

### Examples

```
meroxa apps deploy # assumes you run it from the app directory
meroxa apps deploy --path ./my-app

```

### Options

```
  -h, --help          help for deploy
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

