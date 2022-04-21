## meroxa apps run

Execute a Turbine Data Application locally (Beta)

### Synopsis

meroxa apps run will build your app locally to then run it locally in --path.

```
meroxa apps run [flags]
```

### Examples

```
meroxa apps run 			# assumes you run it from the app directory
meroxa apps run --path ../go-demo 	# it'll use lang defined in your app.json

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

* [meroxa apps](meroxa_apps.md)	 - Manage Turbine Data Applications (Beta)

