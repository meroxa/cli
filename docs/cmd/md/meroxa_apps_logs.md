## meroxa apps logs

View relevant logs to the state of the given Turbine Data Application

### Synopsis

This command will fetch relevant logs about the Application specified in '--path'
(or current working directory if not specified) on our Meroxa Platform,
or the Application specified by the given name or UUID identifier.

```
meroxa apps logs [NameOrUUID] [--path pwd] [flags]
```

### Examples

```
meroxa apps logs # assumes that the Application is in the current directory
meroxa apps logs --path /my/app
meroxa apps logs my-turbine-application
```

### Options

```
  -h, --help          help for logs
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

* [meroxa apps](meroxa_apps.md)	 - Manage Turbine Data Applications

