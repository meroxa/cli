## meroxa apps remove

Remove a Turbine Data Application

### Synopsis

This command will remove the Application specified in '--path'
(or current working directory if not specified) previously deployed on our Meroxa Platform,
or the Application specified by the given name or UUID identifier.

```
meroxa apps remove [NameOrUUID] [--path pwd] [flags]
```

### Examples

```
meroxa apps remove # assumes that the Application is in the current directory
meroxa apps remove --path /my/app
meroxa apps remove NAME
```

### Options

```
  -f, --force         skip confirmation
  -h, --help          help for remove
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

* [meroxa apps](meroxa_apps.md)	 - Manage Turbine Data Applications (Beta)

