## meroxa apps deploy

Deploy a Turbine Data Application

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
      --path string   
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa apps](meroxa_apps.md)	 - Manage Turbine Data Applications

