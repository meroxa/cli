## meroxa apps open

Open the link to a Turbine Data Application in the Dashboard

```
meroxa apps open [--path pwd] [flags]
```

### Examples

```
meroxa apps open # assumes that the Application is in the current directory
meroxa apps open --path /my/app
meroxa apps open NAMEorUUID
```

### Options

```
  -h, --help          help for open
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

