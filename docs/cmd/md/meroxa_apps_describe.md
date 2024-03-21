## meroxa apps describe

Describe a Conduit Data Application

### Synopsis

This command will fetch details about the Application specified in '--path'
(or current working directory if not specified) on our Meroxa Platform,
or the Application specified by the given ID or Application Name.

```
meroxa apps describe [nameOrUUID] [--path pwd] [flags]
```

### Examples

```
meroxa apps describe # assumes that the Application is in the current directory
meroxa apps describe --path /my/app
meroxa apps describe ID
meroxa apps describe NAME 
```

### Options

```
  -h, --help          help for describe
      --path string   Path to the app directory (default is local directory)
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

* [meroxa apps](meroxa_apps.md)	 - Manage Conduit Data Applications

