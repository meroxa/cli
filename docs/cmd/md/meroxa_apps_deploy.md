## meroxa apps deploy

Deploy a Turbine Data Application

### Synopsis

This command will deploy the application specified in '--path'
(or current working directory if not specified) to our Meroxa Platform.
If deployment was successful, you should expect an application you'll be able to fully manage


```
meroxa apps deploy [--path pwd] [flags]
```

### Examples

```
meroxa apps deploy # assumes you run it from the app directory
meroxa apps deploy --path ./my-app

```

### Options

```
      --env string                   environment (name or UUID) where application will be deployed to
  -h, --help                         help for deploy
      --path string                  Path to the app directory (default is local directory)
      --skip-collection-validation   Skips unique destination collection and looping validations
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

