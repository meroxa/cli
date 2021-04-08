## meroxa connect

Connect two resources together

### Synopsis

Use the connect command to automatically configure the connectors required to pull data from one resource 
(source) to another (destination).

This command is equivalent to creating two connectors separately, one from the source to Meroxa and another from Meroxa 
to the destination:

meroxa connect --from RESOURCE-NAME --to RESOURCE-NAME --input SOURCE-INPUT

or

meroxa create connector --from postgres --input accounts # Creates source connector
meroxa create connector --to redshift --input orders # Creates destination connector


```
meroxa connect --from RESOURCE-NAME --to RESOURCE-NAME [flags]
```

### Options

```
      --from string       source resource name
  -h, --help              help for connect
      --input string      command delimeted list of input streams
      --pipeline string   pipeline name to attach connectors to
      --to string         destination resource name
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/meroxa.env)
      --debug           display any debugging information
      --json            output json
```

### SEE ALSO

* [meroxa](meroxa.md)	 - The Meroxa CLI

