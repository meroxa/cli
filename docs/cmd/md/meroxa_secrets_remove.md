## meroxa secrets remove

Remove a Conduit Secret

### Synopsis

This command will remove the secret specified either by name or id

```
meroxa secrets remove [--path pwd] [flags]
```

### Examples

```
meroxa apps remove nameOrUUID
```

### Options

```
  -f, --force   skip confirmation
  -h, --help    help for remove
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

* [meroxa secrets](meroxa_secrets.md)	 - Manage Conduit Data Applications

