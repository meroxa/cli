## meroxa connectors create

Create a connector

### Synopsis

Use `connectors create` to create a connector from a source (--from) or to a destination (--to)

```
meroxa connectors create [NAME] [flags]
```

### Examples

```

meroxa connectors create [NAME] --from pg2kafka --input accounts 
meroxa connectors create [NAME] --to pg2redshift --input orders # --input will be the desired stream 
meroxa connectors create [NAME] --to pg2redshift --input orders --pipeline my-pipeline

```

### Options

```
      --from string       resource name to use as source
  -h, --help              help for create
      --input string      command delimited list of input streams
      --pipeline string   pipeline name to attach connector to
      --to string         resource name to use as destination
```

### Options inherited from parent commands

```
      --config string      config file
      --debug              display any debugging information
      --json               output json
      --timeout duration   set the duration of the client timeout in seconds (default 10s) (default 10s)
```

### SEE ALSO

* [meroxa connectors](meroxa_connectors.md)	 - Manage connectors on Meroxa

