## meroxa config set

Update your Meroxa CLI configuration file with a specific key=value

### Synopsis

This command will let you update your Meroxa configuration file to customize your CLI experience. You can check the presence of this file by running `meroxa config describe`, or even provide your own using `--config my-other-cfg-file`. A key with a format such as `MyKey` will be converted automatically to as `MY_KEY`.

```
meroxa config set [flags]
```

### Examples

```
meroxa config set DisableUpdateNotification=true
meroxa config set DISABLE_UPDATE_NOTIFICATION=true
meroxa config set OneKey=true AnotherKey=false
meroxa config set ApiUrl=https://staging.meroxa.com
```

### Options

```
  -h, --help   help for set
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

* [meroxa config](meroxa_config.md)	 - Manage your Meroxa CLI configuration

