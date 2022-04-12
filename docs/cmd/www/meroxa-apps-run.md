---
createdAt: 
updatedAt: 
title: "meroxa apps run"
slug: meroxa-apps-run
url: /cli/cmd/meroxa-apps-run/
---
## meroxa apps run

Execute a Meroxa Data Application locally

### Synopsis

meroxa apps run will build your app locally to then run it
locally on --path.

```
meroxa apps run [flags]
```

### Examples

```
meroxa apps run # assumes you run it from the app directory
meroxa apps run --path ../js-demo --lang js # in case you didn't specify lang on your app.jsonmeroxa apps run --path ../go-demo # it'll use lang defined in your app.json
```

### Options

```
  -h, --help          help for run
  -l, --lang string   language to use (go | js)
      --path string   path of application to run
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa apps](/cli/cmd/meroxa-apps/)	 - Manage Meroxa Data Applications

