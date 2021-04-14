## meroxa create connector

Create a connector

### Synopsis

Use create connector to create a connector from a source (--from) or to a destination (--to)

```
meroxa create connector [NAME] [flags]
```

### Examples

```

meroxa create connector [NAME] --from pg2kafka --input accounts 
meroxa create connector [NAME] --to pg2redshift --input orders # --input will be the desired stream 
meroxa create connector [NAME] --to pg2redshift --input orders --pipeline my-pipeline

```

### Options

```
      --from string       resource name to use as source
  -h, --help              help for connector
      --input string      command delimited list of input streams
      --pipeline string   pipeline name to attach connector to
      --to string         resource name to use as destination
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/meroxa.env)
      --debug           display any debugging information
      --json            output json
```

### SEE ALSO

* [meroxa create](meroxa_create.md)	 - Create Meroxa pipeline components

